package cache

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	finalizedBlockNumber atomic.Uint64
	latestBlockNumber    atomic.Uint64
)

type toCacheMethod string

const (
	/* legacy interface */
	EthChainID toCacheMethod = "eth_chainId"
	NetVersion toCacheMethod = "net_version"
	/* query based by block hash */
	EthGetBlockByHash                    toCacheMethod = "eth_getBlockByHash"
	EthGetBlockTransactionCountByHash    toCacheMethod = "eth_getBlockTransactionCountByHash"
	EthGetTransactionByBlockHashAndIndex toCacheMethod = "eth_getTransactionByBlockHashAndIndex"

	/* param determine hash */
	EthGetBlockByNumber                    toCacheMethod = "eth_getBlockByNumber"
	EthGetBlockTransactionCountByNumber    toCacheMethod = "eth_getBlockTransactionCountByNumber"
	EthGetTransactionByBlockNumberAndIndex toCacheMethod = "eth_getTransactionByBlockNumberAndIndex"
	EthGetBlockReceipts                    toCacheMethod = "eth_getBlockReceipts"
	EthGetBalance                          toCacheMethod = "eth_getBalance"
	EthGetTransactionCount                 toCacheMethod = "eth_getTransactionCount"
	EthGetCode                             toCacheMethod = "eth_getCode"
	EthCall                                toCacheMethod = "eth_call"
	EthGetStorageAt                        toCacheMethod = "eth_getStorageAt"

	/* short ttl */
	EthBlockNumber          toCacheMethod = "eth_blockNumber"
	EthGasPrice             toCacheMethod = "eth_gasPrice"
	EthMaxPriorityFeePerGas toCacheMethod = "eth_maxPriorityFeePerGas"
)

type blockParamVal string

const (
	Earliest  blockParamVal = "earliest"  // never change
	Finalized blockParamVal = "finalized" // short to medium, 10-20s
	Safe      blockParamVal = "safe"      // short to medium, 10-20s
	Latest    blockParamVal = "latest"    // short 1-2s
	Pending   blockParamVal = "pending"   // no cache
)

var ttlByBlockParamVal = map[blockParamVal]time.Duration{
	Earliest:  longCacheTTL,
	Finalized: longCacheTTL,
	Safe:      mediumCacheTTL,
	Latest:    shortCacheTTL,
}

var blockParamIdxByMethod = map[toCacheMethod]int{
	EthGetBlockByNumber:                    0,
	EthGetBlockTransactionCountByNumber:    0,
	EthGetTransactionByBlockNumberAndIndex: 0,
	EthGetBlockReceipts:                    0,
	EthGetBalance:                          1,
	EthGetTransactionCount:                 1,
	EthGetCode:                             1,
	EthCall:                                1,
	EthGetStorageAt:                        2,
}

const (
	longCacheTTL   = time.Minute * 5
	mediumCacheTTL = time.Second * 10
	shortCacheTTL  = time.Second * 1
)

var cacheTTLByMethod = map[toCacheMethod]time.Duration{
	EthChainID:                           longCacheTTL,
	NetVersion:                           longCacheTTL,
	EthGetBlockByHash:                    longCacheTTL,
	EthGetBlockTransactionCountByHash:    longCacheTTL,
	EthGetTransactionByBlockHashAndIndex: longCacheTTL,
	EthBlockNumber:                       shortCacheTTL,
	EthGasPrice:                          shortCacheTTL,
	EthMaxPriorityFeePerGas:              shortCacheTTL,
}

/*
["0x122A...586d", "latest"]                         // 写法 1：标签
["0x122A...586d", "0x3a8b7a"]                       // 写法 2：区块号
["0x122A...586d", {"blockHash": "0xa1d5...a317"}]   // 写法 3：对象，按区块哈希 -> can cache for longer time
["0x122A...586d", {"blockNumber": "0x3a8b7a"}]      // 写法 4：对象，按区块number
*/
var blockParamOptionalMethods = map[toCacheMethod]struct{}{
	EthGetBalance:          {},
	EthGetTransactionCount: {},
	EthGetCode:             {},
	EthCall:                {},
	EthGetStorageAt:        {},
}

type blockRefKind int

const (
	refNone          blockRefKind = iota
	refTag                        // earliest, finalized, safe, latest, pending
	refBlockNum                   // explicit hex block number
	refBlockHash                  // {blockHash: "0x..."}
	refDefaultLatest              // block param omitted, defaults to latest
)

type blockRef struct {
	kind blockRefKind
	tag  blockParamVal
	num  uint64
}

func resolveBlockRef(method toCacheMethod, params []json.RawMessage) blockRef {
	paramIdx := blockParamIdxByMethod[method]

	if len(params) < paramIdx {
		return blockRef{kind: refNone}
	}

	if len(params) == paramIdx {
		if _, ok := blockParamOptionalMethods[method]; ok {
			return blockRef{kind: refDefaultLatest}
		}
		return blockRef{kind: refNone}
	}

	raw := params[paramIdx]
	paramVal := bytes.Trim(bytes.TrimSpace(raw), `"`)

	tag := blockParamVal(paramVal)
	switch tag {
	case Earliest, Finalized, Safe, Latest, Pending:
		return blockRef{kind: refTag, tag: tag}
	}

	s := string(paramVal)
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		if n, ok := parseBlockNumber(s); ok {
			return blockRef{kind: refBlockNum, num: n}
		}
		return blockRef{kind: refNone}
	}

	if bytes.HasPrefix(paramVal, []byte("{")) {
		var obj struct {
			BlockNumber string `json:"blockNumber"`
			BlockHash   string `json:"blockHash"`
		}
		if err := json.Unmarshal(raw, &obj); err != nil {
			return blockRef{kind: refNone}
		}
		if obj.BlockHash != "" {
			return blockRef{kind: refBlockHash}
		}
		if obj.BlockNumber != "" {
			if n, ok := parseBlockNumber(obj.BlockNumber); ok {
				return blockRef{kind: refBlockNum, num: n}
			}
		}
		return blockRef{kind: refNone}
	}

	return blockRef{kind: refNone}
}

// Cache 提供 JSON-RPC 响应的缓存功能，感知区块最终确认状态。
//
// 缓存策略由 ShouldCache 决定，核心逻辑需要你实现。
// Get/Set/InvalidateOnReorg 的 Redis 操作已实现，你只需实现策略判断。
type Cache interface {
	// ShouldCache 根据 JSON-RPC 方法名和请求参数决定是否缓存及缓存时长。
	//
	// 参数：
	// - method: JSON-RPC 方法名，如 "eth_getBlockByNumber"、"eth_sendRawTransaction"
	// - params: 请求参数数组（原始 JSON），用于提取 block tag
	//
	// 返回值：
	// - shouldCache: 是否应该缓存这个请求的响应
	// - ttl: 缓存有效时间
	//
	// 缓存策略指南：
	// 1. 写方法不缓存：eth_sendRawTransaction, eth_sendTransaction
	// 2. finalized/safe block tag：长缓存（5 分钟），因为 finalized 的数据不会因 reorg 改变
	// 3. latest block tag：短缓存（1-2 秒），因为 latest 随时可能变化
	// 4. pending block tag：不缓存
	// 5. 按区块号查询（如 "0x1234"）：长缓存，但需要判断该区块是否已 finalized
	// 6. 无 block tag 的只读方法（如 eth_chainId, net_version）：长缓存
	//
	// 边界情况：
	// - params 为空或格式异常：安全地降级为不缓存
	// - 未知方法名：默认不缓存
	// - block tag 可能是字符串（"latest"）也可能是十六进制数字（"0x1a2b3c"）
	ShouldCache(method string, params []json.RawMessage) (shouldCache bool, ttl time.Duration)

	// Get 从 Redis 获取缓存的响应。
	Get(ctx context.Context, key string) ([]byte, bool, error)

	// Set 将响应写入 Redis 缓存。
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// InvalidateOnReorg 当检测到链重组时，失效从指定区块开始的所有缓存。
	//
	// 参数：
	// - fromBlock: 发生 reorg 的起始区块号，该区块及之后的缓存都需要失效
	//
	// 实现提示：
	// - 缓存 key 中需要包含区块号信息，以便按范围失效
	// - 可以使用 Redis SCAN + DEL，或维护一个按区块号索引的 key 集合
	// - reorg 可能影响 latest 和 safe 的缓存，finalized 之前的不受影响
	// - 失效操作应该是尽力而为，不要因为失效失败而阻塞请求
	InvalidateOnReorg(ctx context.Context, fromBlock uint64) error
}

type redisCache struct {
	client *redis.Client
}

func New(redisURL string) (Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url %s: %w", redisURL, err)
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis url %s: %w", redisURL, err)
	}

	return &redisCache{client: client}, nil
}

func (c *redisCache) ShouldCache(method string, params []json.RawMessage) (bool, time.Duration) {
	methodName := toCacheMethod(method)

	if ttl, ok := cacheTTLByMethod[methodName]; ok {
		return true, ttl
	}

	if _, ok := blockParamIdxByMethod[methodName]; !ok {
		return false, 0
	}

	ref := resolveBlockRef(methodName, params)
	switch ref.kind {
	case refTag:
		if ttl, ok := ttlByBlockParamVal[ref.tag]; ok {
			return true, ttl
		}
		return false, 0
	case refBlockNum:
		return cacheTTLByBlockNum(ref.num)
	case refBlockHash:
		return true, longCacheTTL
	case refDefaultLatest:
		return true, ttlByBlockParamVal[Latest]
	default:
		return false, 0
	}
}

func parseBlockNumber(s string) (uint64, bool) {
	if len(s) < 3 || (s[:2] != "0x" && s[:2] != "0X") {
		return 0, false
	}
	n, err := strconv.ParseUint(s[2:], 16, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

func cacheTTLByBlockNum(blockNumber uint64) (bool, time.Duration) {
	curFinalized := finalizedBlockNumber.Load()
	if curFinalized == 0 {
		return false, 0
	}

	if blockNumber <= curFinalized {
		return true, longCacheTTL
	}

	if blockNumber <= latestBlockNumber.Load() {
		return true, shortCacheTTL
	}

	return false, 0
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if keyNotExist := errors.Is(err, redis.Nil); keyNotExist {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) InvalidateOnReorg(_ context.Context, _ uint64) error {
	return errors.New("not implemented")
}

func canonicalizeParams(params []json.RawMessage) []json.RawMessage {
	out := make([]json.RawMessage, len(params))
	for i, p := range params {
		out[i] = canonicalizeValue(p)
	}
	return out
}

func canonicalizeValue(raw json.RawMessage) json.RawMessage {
	trimmed := bytes.TrimSpace(raw)

	if len(trimmed) == 0 {
		return raw
	}

	switch trimmed[0] {
	case '"':
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return raw
		}
		s = canonicalizeHex(s)
		b, _ := json.Marshal(s)
		return b

	case '{':
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(trimmed, &obj); err != nil {
			return raw
		}
		normalized := make(map[string]json.RawMessage, len(obj))
		for k, v := range obj {
			if byteDataFields[k] {
				normalized[k] = lowerHexRaw(v)
			} else {
				normK := k
				if strings.HasPrefix(k, "0x") || strings.HasPrefix(k, "0X") {
					normK = lowerHex(k)
				}
				normalized[normK] = canonicalizeValue(v)
			}
		}
		b, _ := json.Marshal(normalized)
		return b

	case '[':
		var arr []json.RawMessage
		if err := json.Unmarshal(trimmed, &arr); err != nil {
			return raw
		}
		for i, v := range arr {
			arr[i] = canonicalizeValue(v)
		}
		b, _ := json.Marshal(arr)
		return b

	default:
		b, _ := json.Marshal(json.RawMessage(trimmed))
		return b
	}
}

var byteDataFields = map[string]bool{
	"data": true, "input": true, "code": true,
}

func lowerHexRaw(raw json.RawMessage) json.RawMessage {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] != '"' {
		return raw
	}
	var s string
	if err := json.Unmarshal(trimmed, &s); err != nil {
		return raw
	}
	s = lowerHex(s)
	b, _ := json.Marshal(s)
	return b
}

func lowerHex(s string) string {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return s
	}
	return "0x" + strings.ToLower(s[2:])
}

func canonicalizeHex(s string) string {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return s
	}
	hex := strings.ToLower(s[2:])
	if len(hex) == 40 {
		return "0x" + hex
	}

	hex = strings.TrimLeft(hex, "0")
	if hex == "" {
		return "0x0"
	}
	return "0x" + hex
}

func BuildKey(chainID uint64, method string, params []json.RawMessage) string {
	params = canonicalizeParams(params)
	h := sha256.New()
	for _, param := range params {
		h.Write(param)
		h.Write([]byte{0}) // add delimiter
	}
	bs := h.Sum(nil)
	paramsHashBase64 := base64.RawURLEncoding.EncodeToString(bs)

	key := fmt.Sprintf("gw:rpc:%d:%s:%s", chainID, method, paramsHashBase64)

	return key
}
