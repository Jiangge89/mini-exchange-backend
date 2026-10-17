package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func New(redisURL string) Cache {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		opt = &redis.Options{Addr: "localhost:6379"}
	}
	return &redisCache{client: redis.NewClient(opt)}
}

func (c *redisCache) ShouldCache(_ string, _ []json.RawMessage) (bool, time.Duration) {
	return false, 0 // not implemented
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
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

func BuildKey(method string, params []json.RawMessage) string {
	h := sha256.New()
	h.Write([]byte(method))
	for _, p := range params {
		h.Write(p)
	}
	return fmt.Sprintf("gw:rpc:%x", h.Sum(nil))
}
