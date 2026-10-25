package cache

import (
	"encoding/json"
	"testing"
)

func TestShouldCache_SendRawTransaction_NeverCache(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_sendRawTransaction 是写方法，永远不应缓存
	// 期望 shouldCache == false
}

func TestShouldCache_FinalizedBlockTag_LongTTL(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_getBlockByNumber with params ["finalized", true]
	// 期望 shouldCache == true, ttl >= 5 分钟
}

func TestShouldCache_LatestBlockTag_ShortTTL(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_getBlockByNumber with params ["latest", true]
	// 期望 shouldCache == true, ttl <= 2 秒
}

func TestShouldCache_PendingBlockTag_NoCache(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_getBlockByNumber with params ["pending", true]
	// 期望 shouldCache == false
}

func TestShouldCache_SpecificBlockNumber_LongTTL(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_getBlockByNumber with params ["0x1234", true]
	// 期望 shouldCache == true，因为具体区块号的数据不会变
	// （但未 finalized 的区块可能因 reorg 而改变，如何处理是一个取舍）
}

func TestShouldCache_ChainId_LongTTL(t *testing.T) {
	t.Skip("TODO: implement")
	// eth_chainId 无参数，结果永远不变
	// 期望 shouldCache == true, ttl 较长
}

func TestShouldCache_EmptyParams_SafeDefault(t *testing.T) {
	t.Skip("TODO: implement")
	// 参数为空或 nil
	// 期望不 panic，返回安全的默认值
}

func TestShouldCache_UnknownMethod_NoCache(t *testing.T) {
	t.Skip("TODO: implement")
	// 未知的方法名
	// 期望 shouldCache == false
}

func TestInvalidateOnReorg_ClearsAffectedBlocks(t *testing.T) {
	t.Skip("TODO: implement")
	// 假设缓存了区块 100-110 的数据
	// reorg 发生在区块 105
	// 期望区块 105-110 的缓存被清除，100-104 的保留
}

// Expected keys are fixed test vectors, computed independently from explicitly
// specified canonical JSON tokens using SHA-256(token + NUL for each parameter)
// and unpadded URL-safe base64. No production helper computes the expected value.
// These vectors describe the current normalization contract, including trimming
// leading zeros in top-level hash strings; they do not validate RPC parameters.
func TestBuildKey(t *testing.T) {
	/*
		blockId 的所有形态

		// 1. Tag 字符串
		"earliest"
		"finalized"
		"safe"
		"latest"
		"pending"

		// 2. Hex 区块号 - prefix 0 should be removed
		"0x0"
		"0x1a2b3c"

		// 3. 对象 — 按 blockHash（仅 eth_getBalance/eth_getTransactionCount/eth_getCode/eth_call/eth_getStorageAt/eth_getProof 等 state 方法 EIP-1898）
		{"blockHash": "0xabc...def"}
		{"blockHash": "0xabc...def", "requireCanonical": true}

		// 4. 对象 — 按 blockNumber
		{"blockNumber": "0x1a2b3c"} - prefix 0 should be removed

		hash 参数 - prefix 0 should NOT be removed

		// 固定格式：32 bytes hex
		"0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd"
	*/
	tests := []struct {
		method  string
		name    string
		params  string
		chainID uint64
		wantKey string
	}{
		// no params
		{method: "eth_chainId", name: "nil", params: `null`, chainID: 1, wantKey: "gw:rpc:1:eth_chainId:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_chainId", name: "empty", params: `[]`, chainID: 1, wantKey: "gw:rpc:1:eth_chainId:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_chainId", name: "other_chain", params: `[]`, chainID: 560048, wantKey: "gw:rpc:560048:eth_chainId:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "net_version", name: "nil", params: `null`, chainID: 1, wantKey: "gw:rpc:1:net_version:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "net_version", name: "empty", params: `[]`, chainID: 1, wantKey: "gw:rpc:1:net_version:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "net_version", name: "other_chain", params: `[]`, chainID: 560048, wantKey: "gw:rpc:560048:net_version:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_blockNumber", name: "nil", params: `null`, chainID: 1, wantKey: "gw:rpc:1:eth_blockNumber:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_blockNumber", name: "empty", params: `[]`, chainID: 1, wantKey: "gw:rpc:1:eth_blockNumber:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_blockNumber", name: "other_chain", params: `[]`, chainID: 560048, wantKey: "gw:rpc:560048:eth_blockNumber:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_gasPrice", name: "nil", params: `null`, chainID: 1, wantKey: "gw:rpc:1:eth_gasPrice:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_gasPrice", name: "empty", params: `[]`, chainID: 1, wantKey: "gw:rpc:1:eth_gasPrice:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_gasPrice", name: "other_chain", params: `[]`, chainID: 560048, wantKey: "gw:rpc:560048:eth_gasPrice:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_maxPriorityFeePerGas", name: "nil", params: `null`, chainID: 1, wantKey: "gw:rpc:1:eth_maxPriorityFeePerGas:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_maxPriorityFeePerGas", name: "empty", params: `[]`, chainID: 1, wantKey: "gw:rpc:1:eth_maxPriorityFeePerGas:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{method: "eth_maxPriorityFeePerGas", name: "other_chain", params: `[]`, chainID: 560048, wantKey: "gw:rpc:560048:eth_maxPriorityFeePerGas:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		// eth_getBlockByHash 2 param: [hash, fullTxs]
		{method: "eth_getBlockByHash", name: "hash_1_False", params: `["0x0000000000000000000000000000000000000000000000000000000000000001",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByHash:u5GSLA0lrTraVqnx-dBzVK264nGMXQkfhbaYFGEb8i0"},
		{method: "eth_getBlockByHash", name: "hash_1_True", params: `["0x0000000000000000000000000000000000000000000000000000000000000001",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByHash:m7kOGGiBlVHx-iuysUZzpkOnpdgaUjrY5R3CYXmHq2M"},
		{method: "eth_getBlockByHash", name: "hash_2_False", params: `["0x0000000000000000000000000000000000000000000000000000000000000002",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByHash:Mymyw6qaUp_I7eAKq0LUXduSXFH6ZtUcTfK-WmmOO8A"},
		{method: "eth_getBlockByHash", name: "hash_2_True", params: `["0x0000000000000000000000000000000000000000000000000000000000000002",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByHash:vYRWjZoaOe6xCnOnDI3m4U98l43YW5TsojBE3MEbhrw"},
		// eth_getBlockByNumber 2 params: [blockId, fullTxs]
		{method: "eth_getBlockByNumber", name: "earliest_False", params: `["earliest",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:_SSBVIH9bNMVZ4fUoPch9Qi7BNhq5q7lOce0x3zhTPs"},
		{method: "eth_getBlockByNumber", name: "earliest_True", params: `["earliest",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:yJ7_2tbLXZk1HF6KODzD0J5V5sz8Tqh50boVNH-roKQ"},
		{method: "eth_getBlockByNumber", name: "finalized_False", params: `["finalized",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:7h_4TeM9Sm5wOsrwRLG0O1hgHAB25EqhVYrR6qhBDUo"},
		{method: "eth_getBlockByNumber", name: "finalized_True", params: `["finalized",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:e1Zp4vz9za0b8eqrIwwLQMWm-RBcZuI9m6kIBWuCwRg"},
		{method: "eth_getBlockByNumber", name: "safe_False", params: `["safe",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:s7MWXfFGFMS6uDXUYKcyPLgOxFNG0HOWtDEcEQfTHXk"},
		{method: "eth_getBlockByNumber", name: "safe_True", params: `["safe",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:-kjg2e2xTNMdrOyu-mVJ3816z-nZf9-W76s3jt32oIo"},
		{method: "eth_getBlockByNumber", name: "latest_False", params: `["latest",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:36eyvZaleev25x3EuBsQwVxtKiniUWVzeLkCZbLq90M"},
		{method: "eth_getBlockByNumber", name: "latest_True", params: `["latest",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:BSrmmNp1s3AqTM20xMGSnT-YBsO4v4g-aK97lSUfIOg"},
		{method: "eth_getBlockByNumber", name: "pending_False", params: `["pending",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:t4FhrHBNdfPI0YfXKTDN2kNj5GGBF9ZQAgP9sXZs7VM"},
		{method: "eth_getBlockByNumber", name: "pending_True", params: `["pending",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:LF-Q_lv-XXgt8zsPf0946jvnG6PpFaLA6Y8kQwA94JA"},
		{method: "eth_getBlockByNumber", name: "0x0_False", params: `["0x0",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:AJsjVNyxSXIdw4e0IFMdRd6CGWhCPPYhOhxJDoKi7F8"},
		{method: "eth_getBlockByNumber", name: "0x0_True", params: `["0x0",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:WXCDkM8QzZ5DO6AI3qciqg33e54QY1f6eohjQJgS2lE"},
		{method: "eth_getBlockByNumber", name: "0x1234_False", params: `["0x1234",false]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:lZGQqyL0GlF2sUanVqUhHfgzA4GMISIrv45HobKRtzg"},
		{method: "eth_getBlockByNumber", name: "0x1234_True", params: `["0x1234",true]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockByNumber:0q1MotSTEPETnTH-D8Fzg3AF3LHrEzKRXJk6LMXxlL8"},
		// 1 param:  [hash]
		{method: "eth_getBlockTransactionCountByHash", name: "hash_1_only", params: `["0x0000000000000000000000000000000000000000000000000000000000000001"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByHash:3rjsaibKwpMoFojyCJHonvayA5TJxjG88LRz1v_qHMM"},
		{method: "eth_getBlockTransactionCountByHash", name: "hash_2_only", params: `["0x0000000000000000000000000000000000000000000000000000000000000002"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByHash:b8aQ0x4YUVZ140BGXatwIgCyWwhaGmk_n0Wjmr8aTq4"},
		// 1 param: [blockId]
		{method: "eth_getBlockTransactionCountByNumber", name: "earliest", params: `["earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:OkhMlUZKCfLQ0N4zMnqhn4ZJ5IMt8FnqUZr9mFKDgOI"},
		{method: "eth_getBlockTransactionCountByNumber", name: "finalized", params: `["finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:Wu2oVYwwAVL_OvPUoiuf_qCrFQN6q4_ccPHTZ0FErvk"},
		{method: "eth_getBlockTransactionCountByNumber", name: "safe", params: `["safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:mEO6I74jDwcEwqmOJhuT3EvEUBKA0IyLpUPDBGijvMY"},
		{method: "eth_getBlockTransactionCountByNumber", name: "latest", params: `["latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:ESCQ-1tN2kRIddTQs7HgxQlGcYVNl9Ak_I4AdeseXFQ"},
		{method: "eth_getBlockTransactionCountByNumber", name: "pending", params: `["pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:-I_pozNhAWU4CWBZ4ysqtC_x4IwHD6152NyT5L0DRms"},
		{method: "eth_getBlockTransactionCountByNumber", name: "0x0", params: `["0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:HqokD3TaM5voYNBUX-L4_Cw5t942-p9gAwbJiO3Lm_E"},
		{method: "eth_getBlockTransactionCountByNumber", name: "0x1234", params: `["0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockTransactionCountByNumber:Quk7Ql1rdphqG3UD465lWsd3MDuyzVMssGTpQmpeCuw"},
		// 2 params: [hash, index]
		{method: "eth_getTransactionByBlockHashAndIndex", name: "hash_1_0x0", params: `["0x0000000000000000000000000000000000000000000000000000000000000001","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockHashAndIndex:xP4HjtJoiZ9uh8lfVJsW-seF0HmM604eEQv-Uhg9ANY"},
		{method: "eth_getTransactionByBlockHashAndIndex", name: "hash_1_0x1", params: `["0x0000000000000000000000000000000000000000000000000000000000000001","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockHashAndIndex:_ACstkQKI9q9d2vQaFe-hvbzsU8tECT9XrtvqS1e6YU"},
		{method: "eth_getTransactionByBlockHashAndIndex", name: "hash_2_0x0", params: `["0x0000000000000000000000000000000000000000000000000000000000000002","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockHashAndIndex:P6x9xn5jzH97U2k-GGgpsRMBzbvelOBA2aSe-6z08AA"},
		{method: "eth_getTransactionByBlockHashAndIndex", name: "hash_2_0x1", params: `["0x0000000000000000000000000000000000000000000000000000000000000002","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockHashAndIndex:dmZ6HnxH3UA8g17zvaVk420x9Ay0ilmxsO6NhJ1dln0"},
		// 2 params: [blockId, index]
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "earliest_0x0", params: `["earliest","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:ifUDX_dmw9NLB3VSUjcKIIMAJg_KrNZ-Ar4dqGzA4_c"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "earliest_0x1", params: `["earliest","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:Gs_ZvVLrGGEgrOx-P99ZtbF3G8TmOYtHm6N5H8bZjzE"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "finalized_0x0", params: `["finalized","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:pD7thjGypQ1_X2Nf5qWf2FQnuGR1hawUnjUne17I-4U"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "finalized_0x1", params: `["finalized","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:2O7eY7BuiiKBoZdxMgQsThnFYC7j8b808QVM_f1owuM"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "safe_0x0", params: `["safe","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:mGAFLnLeaeysM7buQfg4tAp59Me2x5T4QXT7Ku2nzIs"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "safe_0x1", params: `["safe","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:YctBmNRVXaHPQHmypEfuyvYsK8PlKUOy5m-ZUha1nyg"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "latest_0x0", params: `["latest","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:Al4cne-EFr_hkjywmLD_SU_trOtfTQTKGGHM6nHm45A"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "latest_0x1", params: `["latest","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:e--NPugm21nZzGGpD3ToB9azietfCP72NSFzWxSDtXY"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "pending_0x0", params: `["pending","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:jBFpsLPXXnXxReFMrL-2sQclF6cgAnEeNsSP3jxknRo"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "pending_0x1", params: `["pending","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:gfUe5nxeruOGUVancUo3BbhFi6fSW9khv-FG-HeYg5A"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "0x0_0x0", params: `["0x0","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:R9mu6C-QzHfebKFOej5AGYZ3J-02rrp3_Eke0CBdTOY"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "0x0_0x1", params: `["0x0","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:0zzJQkwiyfUhtyNOZYtun2wsCGMiHQeYO3avPqXsJ3w"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "0x1234_0x0", params: `["0x1234","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:z_b8P_ZGH8tWoFnD1FeoKqrnxTisWbbmm6zL7S45NCk"},
		{method: "eth_getTransactionByBlockNumberAndIndex", name: "0x1234_0x1", params: `["0x1234","0x1"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionByBlockNumberAndIndex:z766v3orbWWLgeB6ON4IUK3dMWICK4iaKujG-MQNlas"},
		// 1 param: [blockId]
		{method: "eth_getBlockReceipts", name: "earliest", params: `["earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:OkhMlUZKCfLQ0N4zMnqhn4ZJ5IMt8FnqUZr9mFKDgOI"},
		{method: "eth_getBlockReceipts", name: "finalized", params: `["finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:Wu2oVYwwAVL_OvPUoiuf_qCrFQN6q4_ccPHTZ0FErvk"},
		{method: "eth_getBlockReceipts", name: "safe", params: `["safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:mEO6I74jDwcEwqmOJhuT3EvEUBKA0IyLpUPDBGijvMY"},
		{method: "eth_getBlockReceipts", name: "latest", params: `["latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:ESCQ-1tN2kRIddTQs7HgxQlGcYVNl9Ak_I4AdeseXFQ"},
		{method: "eth_getBlockReceipts", name: "pending", params: `["pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:-I_pozNhAWU4CWBZ4ysqtC_x4IwHD6152NyT5L0DRms"},
		{method: "eth_getBlockReceipts", name: "0x0", params: `["0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:HqokD3TaM5voYNBUX-L4_Cw5t942-p9gAwbJiO3Lm_E"},
		{method: "eth_getBlockReceipts", name: "0x1234", params: `["0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:Quk7Ql1rdphqG3UD465lWsd3MDuyzVMssGTpQmpeCuw"},
		{method: "eth_getBlockReceipts", name: "block_number", params: `[{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:UAuMW2BFn79zaN279hBGjGRhKyVpy8hvdaV6oqwybMM"},
		{method: "eth_getBlockReceipts", name: "block_hash", params: `[{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:pGTl4Gxc19SdKV3HaSgN0XBBsS55tiVTOSPhmNnEklM"},
		{method: "eth_getBlockReceipts", name: "canonical_hash", params: `[{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBlockReceipts:Md_fXcBEigs3C3LPifPI2XERc2OOJqFPVYHX1Qd4pnk"},
		// 2 params: [address, blockId?]
		{method: "eth_getBalance", name: "earliest", params: `["0x0000000000000000000000000000000000000001","earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:bmi4-3Zd5LOwCK60dOrUL0h_cdnIPMk5nA9dLZSmaCs"},
		{method: "eth_getBalance", name: "finalized", params: `["0x0000000000000000000000000000000000000001","finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:hSQkRd9YM8oXOutgAE9IDJvfhGa3YRKV38wjnXgu6rE"},
		{method: "eth_getBalance", name: "safe", params: `["0x0000000000000000000000000000000000000001","safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:WwRxIJYzsENXWvzkWqf9VoghggpZYuQRYt-z5N5TkOk"},
		{method: "eth_getBalance", name: "latest", params: `["0x0000000000000000000000000000000000000001","latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:JNKmOsU1bI8EHCmvQKnDs_4yiB-ERTRxFHLapqDuG2s"},
		{method: "eth_getBalance", name: "pending", params: `["0x0000000000000000000000000000000000000001","pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:zmfLYi5Ol750p3MYstx6YXuOOwXE-y5jInYkqjMF4pc"},
		{method: "eth_getBalance", name: "0x0", params: `["0x0000000000000000000000000000000000000001","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:sypO_qGMPkNQMo9okAAQWWYt61LKuZXzX8HSUya_PJI"},
		{method: "eth_getBalance", name: "0x1234", params: `["0x0000000000000000000000000000000000000001","0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:d4TEpYAwX1n4rI6ZSShRKzwo_8pwHC39BM3RGYm1Z9g"},
		{method: "eth_getBalance", name: "block_number", params: `["0x0000000000000000000000000000000000000001",{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:o-gtsgJn55FtJW132PCn4DjSA-UN2IcM55_dZTs8j1M"},
		{method: "eth_getBalance", name: "block_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:7MVqSYouXMTLg3LEqw6Lli8m2JqoHxO8IdRXHalgfjU"},
		{method: "eth_getBalance", name: "canonical_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:ggcSPnxiZlPeP7hcjKEIot9FqzE9EGh7E20vWTMQSxM"},
		{method: "eth_getBalance", name: "omitted_block", params: `["0x0000000000000000000000000000000000000001"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:H_N4QDim5_qwWs79D9f9NeldZFRMVfhCBRAS2WTUC-k"},
		{method: "eth_getBalance", name: "uppercase_address", params: `["0x00000000000000000000000000000000000000AB","latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getBalance:4mo6nFipgSFnSWfRDUieEaJpmayXz6uDRGXQbiPiCYs"},
		//  [address, blockId?]
		{method: "eth_getTransactionCount", name: "earliest", params: `["0x0000000000000000000000000000000000000001","earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:bmi4-3Zd5LOwCK60dOrUL0h_cdnIPMk5nA9dLZSmaCs"},
		{method: "eth_getTransactionCount", name: "finalized", params: `["0x0000000000000000000000000000000000000001","finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:hSQkRd9YM8oXOutgAE9IDJvfhGa3YRKV38wjnXgu6rE"},
		{method: "eth_getTransactionCount", name: "safe", params: `["0x0000000000000000000000000000000000000001","safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:WwRxIJYzsENXWvzkWqf9VoghggpZYuQRYt-z5N5TkOk"},
		{method: "eth_getTransactionCount", name: "latest", params: `["0x0000000000000000000000000000000000000001","latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:JNKmOsU1bI8EHCmvQKnDs_4yiB-ERTRxFHLapqDuG2s"},
		{method: "eth_getTransactionCount", name: "pending", params: `["0x0000000000000000000000000000000000000001","pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:zmfLYi5Ol750p3MYstx6YXuOOwXE-y5jInYkqjMF4pc"},
		{method: "eth_getTransactionCount", name: "0x0", params: `["0x0000000000000000000000000000000000000001","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:sypO_qGMPkNQMo9okAAQWWYt61LKuZXzX8HSUya_PJI"},
		{method: "eth_getTransactionCount", name: "0x1234", params: `["0x0000000000000000000000000000000000000001","0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:d4TEpYAwX1n4rI6ZSShRKzwo_8pwHC39BM3RGYm1Z9g"},
		{method: "eth_getTransactionCount", name: "block_number", params: `["0x0000000000000000000000000000000000000001",{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:o-gtsgJn55FtJW132PCn4DjSA-UN2IcM55_dZTs8j1M"},
		{method: "eth_getTransactionCount", name: "block_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:7MVqSYouXMTLg3LEqw6Lli8m2JqoHxO8IdRXHalgfjU"},
		{method: "eth_getTransactionCount", name: "canonical_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:ggcSPnxiZlPeP7hcjKEIot9FqzE9EGh7E20vWTMQSxM"},
		{method: "eth_getTransactionCount", name: "omitted_block", params: `["0x0000000000000000000000000000000000000001"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getTransactionCount:H_N4QDim5_qwWs79D9f9NeldZFRMVfhCBRAS2WTUC-k"},
		// [address, blockId?]
		{method: "eth_getCode", name: "earliest", params: `["0x0000000000000000000000000000000000000001","earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:bmi4-3Zd5LOwCK60dOrUL0h_cdnIPMk5nA9dLZSmaCs"},
		{method: "eth_getCode", name: "finalized", params: `["0x0000000000000000000000000000000000000001","finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:hSQkRd9YM8oXOutgAE9IDJvfhGa3YRKV38wjnXgu6rE"},
		{method: "eth_getCode", name: "safe", params: `["0x0000000000000000000000000000000000000001","safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:WwRxIJYzsENXWvzkWqf9VoghggpZYuQRYt-z5N5TkOk"},
		{method: "eth_getCode", name: "latest", params: `["0x0000000000000000000000000000000000000001","latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:JNKmOsU1bI8EHCmvQKnDs_4yiB-ERTRxFHLapqDuG2s"},
		{method: "eth_getCode", name: "pending", params: `["0x0000000000000000000000000000000000000001","pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:zmfLYi5Ol750p3MYstx6YXuOOwXE-y5jInYkqjMF4pc"},
		{method: "eth_getCode", name: "0x0", params: `["0x0000000000000000000000000000000000000001","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:sypO_qGMPkNQMo9okAAQWWYt61LKuZXzX8HSUya_PJI"},
		{method: "eth_getCode", name: "0x1234", params: `["0x0000000000000000000000000000000000000001","0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:d4TEpYAwX1n4rI6ZSShRKzwo_8pwHC39BM3RGYm1Z9g"},
		{method: "eth_getCode", name: "block_number", params: `["0x0000000000000000000000000000000000000001",{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:o-gtsgJn55FtJW132PCn4DjSA-UN2IcM55_dZTs8j1M"},
		{method: "eth_getCode", name: "block_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:7MVqSYouXMTLg3LEqw6Lli8m2JqoHxO8IdRXHalgfjU"},
		{method: "eth_getCode", name: "canonical_hash", params: `["0x0000000000000000000000000000000000000001",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:ggcSPnxiZlPeP7hcjKEIot9FqzE9EGh7E20vWTMQSxM"},
		{method: "eth_getCode", name: "omitted_block", params: `["0x0000000000000000000000000000000000000001"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getCode:H_N4QDim5_qwWs79D9f9NeldZFRMVfhCBRAS2WTUC-k"},
		//  [txObject, blockId?, stateOverrides?]
		{method: "eth_call", name: "earliest", params: `[{},"earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:g1y1edUvvyj53cDF8uQw0AJKrRADBjuoaUoMH39XT1A"},
		{method: "eth_call", name: "finalized", params: `[{},"finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:10GEGRP6pF181PNJJw_yds52VPRJ-9HLmRGJej0vlNk"},
		{method: "eth_call", name: "safe", params: `[{},"safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:UpZ6RYbcgX8Hoh_1jxXAYP26dYEjQ2uwNplaGNiGC3I"},
		{method: "eth_call", name: "latest", params: `[{},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:nXl7Y18xXUE0NsKlB01gHdQlLEbqSfh29mKjb64wFrU"},
		{method: "eth_call", name: "pending", params: `[{},"pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:H45TjSOfYsqKmGBMW6S2ZAfnu8gQT3uT7aMNi2nfg48"},
		{method: "eth_call", name: "0x0", params: `[{},"0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:KdmsHZBoydh_ghDsMWqc-OMn0NH1FMuJ0fZYMoHUJEE"},
		{method: "eth_call", name: "0x1234", params: `[{},"0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:nftjpRYn2we2AXPv6eQ3nRT5e7d5bVQjKFM9Lvg50xk"},
		{method: "eth_call", name: "block_number", params: `[{},{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:FoF-K49mYkqaRJQxhj7qjNuh_xJSO-w0Gg-k1brb7OQ"},
		{method: "eth_call", name: "block_hash", params: `[{},{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:EQlSjfc8eIHUw1GQLw3a77VF9VatAoAp8tt_5ssgqwk"},
		{method: "eth_call", name: "canonical_hash", params: `[{},{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:v8n9RCfAAvTKApUTIMH_EPp9E_0L5hxhm1OSVkFjyKg"},
		{method: "eth_call", name: "omitted_block", params: `[{}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:aOnoa2kmzCs335a15hu4yr2rJ2Jyu3P2dn5CDD6tBmM"},
		{method: "eth_call", name: "to_0", params: `[{"to":"0x0000000000000000000000000000000000000001"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:JvTqKjUFjapvcDwUibsjSQn8_hXeHEexXwy7SmnH8vk"},
		{method: "eth_call", name: "to_1", params: `[{"to":"0x0000000000000000000000000000000000000002"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:JS-N1ZpHGS1Sd2kj74bSyv-EYNtYNExnlIuUPR5iKEc"},
		{method: "eth_call", name: "from_0", params: `[{"from":"0x0000000000000000000000000000000000000001"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:EMi2zTP8Ux9apOJOUYJ_YjeeehkzGHyO8GAiA6lJMQI"},
		{method: "eth_call", name: "data_0", params: `[{"data":"0x"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:EZz-PROtc7kT9-OT_5k6mr7Rq-l6fsM45m3h9mdGvx8"},
		{method: "eth_call", name: "data_1", params: `[{"data":"0x00"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:D0rc232nI7dlL8nkkYoLH4V9IC5ddFq4K0Iy9XrP900"},
		{method: "eth_call", name: "data_2", params: `[{"data":"0x01"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:viRLZI4VSe9-mGghJe2tb7ZBU_LtaTYtsixVoPX9L3Y"},
		{method: "eth_call", name: "data_3", params: `[{"data":"0x0001"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:76ICeamP0M7IX9sHpjBN5_BbMlMf-7_X4H58xsWpumU"},
		{method: "eth_call", name: "input_0", params: `[{"input":"0x0001"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:5gHwOuTQI-csLocv65r_FH__8lZWDAWc41gLJ6v1DaU"},
		{method: "eth_call", name: "gas_0", params: `[{"gas":"0x5208"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:X58E0cAiz8_tFxC9SJo8qFh5ztw27iOfl2NE7O8RQb0"},
		{method: "eth_call", name: "gasPrice_0", params: `[{"gasPrice":"0x1"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:cSry0sOUU6kbqMd-Nn8ewEcWFhWDT0SAYS-dTyNTnHg"},
		{method: "eth_call", name: "value_0", params: `[{"value":"0x0"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:UgK3nd0ll5qsnwAZGPWpm0kQROCwDRlfFdDlVaBA1Eo"},
		{method: "eth_call", name: "value_1", params: `[{"value":"0x1"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:BTVveYPgyHKxxdNWw1M0fkFofR012CHCqq_U34RungE"},
		{method: "eth_call", name: "maxFeePerGas_0", params: `[{"maxFeePerGas":"0x2"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:UBp1yGatktPtt-EIlZPDXoMzOfpHUHNfk57IjXA_XXU"},
		{method: "eth_call", name: "maxPriorityFeePerGas_0", params: `[{"maxPriorityFeePerGas":"0x1"},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:srypq4IGXNSr7sY4A6hzIk__6-SDK9RFDvPXIz-Xhf8"},
		{method: "eth_call", name: "accessList_0", params: `[{"accessList":[]},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:ljnl2BcZPt9mP2YSJD4JB7iq67akz0TYUKR3neQaW9Q"},
		{method: "eth_call", name: "accessList_1entry", params: `[{"accessList":[{"address":"0x0000000000000000000000000000000000000001","storageKeys":["0x0000000000000000000000000000000000000000000000000000000000000001"]}]},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:lhMPyKpOb8zUxZFMfFPaRDSa1FeYIabaHDfHOrwaxXE"},
		{method: "eth_call", name: "accessList_2keys", params: `[{"accessList":[{"address":"0x0000000000000000000000000000000000000001","storageKeys":["0x0000000000000000000000000000000000000000000000000000000000000001","0x0000000000000000000000000000000000000000000000000000000000000002"]}]},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:BqFRnRFL5xn-PuDpzLevHGT5jSJHvaw27c8pLAVkRMo"},
		{method: "eth_call", name: "accessList_2entries", params: `[{"accessList":[{"address":"0x0000000000000000000000000000000000000001","storageKeys":[]},{"address":"0x0000000000000000000000000000000000000002","storageKeys":[]}]},"latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:ppexQ79-KgpC9YFS50bUe1btbBJejo9MXdU3WMIgu4E"},
		{method: "eth_call", name: "override_empty", params: `[{},"latest",{}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:5NclP0oYsuuvHL6iRvJdO5OYEvQkWgo6WM8oujY3FEs"},
		{method: "eth_call", name: "override_balance", params: `[{},"latest",{"0x0000000000000000000000000000000000000001":{"balance":"0x1"}}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:MgbaH-FC-C0VJUDac-W9cOWqJfTVaS2T0eXzf_NnLLg"},
		{method: "eth_call", name: "override_nonce", params: `[{},"latest",{"0x0000000000000000000000000000000000000001":{"nonce":"0x1"}}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:sFFe6S8au6DZGMtiRFXw9Ade7O5KjR8EcDXeqz4rk2M"},
		{method: "eth_call", name: "override_code", params: `[{},"latest",{"0x0000000000000000000000000000000000000001":{"code":"0x0001"}}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:teJ9LWfoS91YH9AsWw9ze2fiFPErEdhTy6klpkYfrN0"},
		{method: "eth_call", name: "override_state", params: `[{},"latest",{"0x0000000000000000000000000000000000000001":{"state":{"0x0000000000000000000000000000000000000000000000000000000000000001":"0x0000000000000000000000000000000000000000000000000000000000000001"}}}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:1GEy-12WHwGfGDz9CxLTRXMvqB_VU_k7keEnyxiO-MA"},
		{method: "eth_call", name: "override_state_diff", params: `[{},"latest",{"0x0000000000000000000000000000000000000001":{"stateDiff":{"0x0000000000000000000000000000000000000000000000000000000000000001":"0x0000000000000000000000000000000000000000000000000000000000000001"}}}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:Pho4kdaU3x_1NMy7pTls8wR6iB5kbpyz8aUeQTZ_MyA"},
		{method: "eth_call", name: "escaped_tag", params: `[{},"lat\u0065st"]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:nXl7Y18xXUE0NsKlB01gHdQlLEbqSfh29mKjb64wFrU"},
		{method: "eth_call", name: "object_block_number_leading_zeros", params: `[{}, {"blockNumber":"0x001a"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_call:MqyKz8Gtb9vhVNhhP182nLJ3tjgQijqmC_v1Ma-PEKU"},
		// [address, slot, blockId?]
		{method: "eth_getStorageAt", name: "earliest", params: `["0x0000000000000000000000000000000000000001","0x0","earliest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:xB5tJ-Uu7WjQ_Jm-QgWx2tszUNJEZUTArFzhp-EOjXw"},
		{method: "eth_getStorageAt", name: "finalized", params: `["0x0000000000000000000000000000000000000001","0x0","finalized"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:9AwKD_nbqNEveK-oq6OuPx5FCqYLmlhb5brQeXjoEj4"},
		{method: "eth_getStorageAt", name: "safe", params: `["0x0000000000000000000000000000000000000001","0x0","safe"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:3wpiWrqLvCHbeoyp6jM3Pq7dfFkN-W5hYNvjAyQTMB8"},
		{method: "eth_getStorageAt", name: "latest", params: `["0x0000000000000000000000000000000000000001","0x0","latest"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:TLsyW4siJMcI1AjNk1FnJnaEO8UltPdswkfPQ9W7goU"},
		{method: "eth_getStorageAt", name: "pending", params: `["0x0000000000000000000000000000000000000001","0x0","pending"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:z1aTT27lEqPvOe1cEV6aT9hK9vla5Mo_fQZYnsRja4c"},
		{method: "eth_getStorageAt", name: "0x0", params: `["0x0000000000000000000000000000000000000001","0x0","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:b5AEAXTmDaPCYHoJzYWv36D2IC3wc5E-sATzCgcy_5s"},
		{method: "eth_getStorageAt", name: "0x1234", params: `["0x0000000000000000000000000000000000000001","0x0","0x1234"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:XTvNYIwxcrX2MtaTZDkI5nGbI3axCilzpqrYPrbnJho"},
		{method: "eth_getStorageAt", name: "block_number", params: `["0x0000000000000000000000000000000000000001","0x0",{"blockNumber":"0x1234"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:dRDfnOObnHxtHRtxVE9GVBeEjuSr0hvpygfHUXMZpW4"},
		{method: "eth_getStorageAt", name: "block_hash", params: `["0x0000000000000000000000000000000000000001","0x0",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:dZznS0IJSbwwuEX-t9nOsRkgJxIDM4DN8ypZfoXbn-U"},
		{method: "eth_getStorageAt", name: "canonical_hash", params: `["0x0000000000000000000000000000000000000001","0x0",{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000001","requireCanonical":true}]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:qFoMCRpUe8g_l3YtkrQLbf_OApit_Ys-UGhPF3IhCkY"},
		{method: "eth_getStorageAt", name: "omitted_block", params: `["0x0000000000000000000000000000000000000001","0x0"]`, chainID: 1, wantKey: "gw:rpc:1:eth_getStorageAt:sypO_qGMPkNQMo9okAAQWWYt61LKuZXzX8HSUya_PJI"},
	}
	for _, tt := range tests {
		t.Run(tt.method+"/"+tt.name, func(t *testing.T) {
			var params []json.RawMessage
			if err := json.Unmarshal([]byte(tt.params), &params); err != nil {
				t.Fatalf("invalid test params: %v", err)
			}
			if got := BuildKey(tt.chainID, tt.method, params); got != tt.wantKey {
				t.Fatalf("BuildKey(%d, %q, %s) = %q; want %q", tt.chainID, tt.method, tt.params, got, tt.wantKey)
			}
		})
	}
}

func mustParseParams(t *testing.T, s string) []json.RawMessage {
	t.Helper()
	var params []json.RawMessage
	if err := json.Unmarshal([]byte(s), &params); err != nil {
		t.Fatalf("invalid params JSON: %v", err)
	}
	return params
}

func TestBuildKey_LeadingZeroNormalization(t *testing.T) {
	const addr = "0x0000000000000000000000000000000000000001"
	const hash = "0x0000000000000000000000000000000000000000000000000000000000000001"

	tests := []struct {
		method string
		name   string
		a, b   string
	}{
		// eth_getBlockByNumber [blockId, fullTxs]
		{"eth_getBlockByNumber", "blockId", `["0x1234",false]`, `["0x01234",false]`},
		// eth_getBlockTransactionCountByNumber [blockId]
		{"eth_getBlockTransactionCountByNumber", "blockId", `["0x1234"]`, `["0x01234"]`},
		// eth_getTransactionByBlockNumberAndIndex [blockId, index]
		{"eth_getTransactionByBlockNumberAndIndex", "blockId", `["0x1234","0x0"]`, `["0x01234","0x0"]`},
		{"eth_getTransactionByBlockNumberAndIndex", "index", `["0x1234","0x1"]`, `["0x1234","0x01"]`},
		// eth_getTransactionByBlockHashAndIndex [hash, index]
		{"eth_getTransactionByBlockHashAndIndex", "index", `["` + hash + `","0x0"]`, `["` + hash + `","0x00"]`},
		// eth_getBlockReceipts [blockId]
		{"eth_getBlockReceipts", "blockId", `["0x1234"]`, `["0x01234"]`},
		// eth_getBalance [address, blockId]
		{"eth_getBalance", "blockId", `["` + addr + `","0x1234"]`, `["` + addr + `","0x01234"]`},
		{"eth_getBalance", "obj_blockNumber", `["` + addr + `",{"blockNumber":"0x1234"}]`, `["` + addr + `",{"blockNumber":"0x01234"}]`},
		// eth_getTransactionCount [address, blockId]
		{"eth_getTransactionCount", "blockId", `["` + addr + `","0x1234"]`, `["` + addr + `","0x01234"]`},
		// eth_getCode [address, blockId]
		{"eth_getCode", "blockId", `["` + addr + `","0x1234"]`, `["` + addr + `","0x01234"]`},
		// eth_call [txObject, blockId]
		{"eth_call", "blockId", `[{},"0x1234"]`, `[{},"0x01234"]`},
		{"eth_call", "obj_blockNumber", `[{},{"blockNumber":"0x1a"}]`, `[{},{"blockNumber":"0x001a"}]`},
		{"eth_call", "txObject_gas", `[{"gas":"0x5208"},"latest"]`, `[{"gas":"0x05208"},"latest"]`},
		{"eth_call", "txObject_value", `[{"value":"0x1"},"latest"]`, `[{"value":"0x01"},"latest"]`},
		// eth_getStorageAt [address, slot, blockId]
		{"eth_getStorageAt", "blockId", `["` + addr + `","0x0","0x1234"]`, `["` + addr + `","0x0","0x01234"]`},
		{"eth_getStorageAt", "slot", `["` + addr + `","0x0","latest"]`, `["` + addr + `","0x00","latest"]`},
		{"eth_getStorageAt", "obj_blockNumber", `["` + addr + `","0x0",{"blockNumber":"0x1234"}]`, `["` + addr + `","0x0",{"blockNumber":"0x01234"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.method+"/"+tt.name, func(t *testing.T) {
			keyA := BuildKey(1, tt.method, mustParseParams(t, tt.a))
			keyB := BuildKey(1, tt.method, mustParseParams(t, tt.b))
			if keyA != keyB {
				t.Fatalf("leading-zero normalization failed:\n  %s → %s\n  %s → %s", tt.a, keyA, tt.b, keyB)
			}
		})
	}
}

func TestBuildKey_CaseNormalization(t *testing.T) {
	const addrLower = "0x000000000000000000000000000000000000abcd"
	const addrUpper = "0x000000000000000000000000000000000000ABCD"
	const hashLower = "0x000000000000000000000000000000000000000000000000000000000000abcd"
	const hashUpper = "0x000000000000000000000000000000000000000000000000000000000000ABCD"

	tests := []struct {
		method string
		name   string
		a, b   string
	}{
		// block number case: "0x1a" vs "0x1A"
		{"eth_getBlockByNumber", "blockId", `["0x1a",false]`, `["0x1A",false]`},
		{"eth_getBlockByNumber", "0x_prefix", `["0x1a",false]`, `["0X1a",false]`},
		{"eth_getBlockTransactionCountByNumber", "blockId", `["0x1a"]`, `["0x1A"]`},
		{"eth_getTransactionByBlockNumberAndIndex", "blockId", `["0x1a","0x0"]`, `["0x1A","0x0"]`},
		{"eth_getTransactionByBlockNumberAndIndex", "index", `["0x1a","0xa"]`, `["0x1a","0xA"]`},
		{"eth_getBlockReceipts", "blockId", `["0x1a"]`, `["0x1A"]`},
		// hash case: 64-char hex, zeros preserved, only case changes
		{"eth_getBlockByHash", "hash", `["` + hashLower + `",false]`, `["` + hashUpper + `",false]`},
		{"eth_getBlockTransactionCountByHash", "hash", `["` + hashLower + `"]`, `["` + hashUpper + `"]`},
		{"eth_getTransactionByBlockHashAndIndex", "hash", `["` + hashLower + `","0x0"]`, `["` + hashUpper + `","0x0"]`},
		// address case: 40-char hex, zeros preserved, only case changes
		{"eth_getBalance", "address", `["` + addrLower + `","latest"]`, `["` + addrUpper + `","latest"]`},
		{"eth_getBalance", "blockId", `["` + addrLower + `","0x1a"]`, `["` + addrLower + `","0x1A"]`},
		{"eth_getBalance", "obj_blockHash", `["` + addrLower + `",{"blockHash":"` + hashLower + `"}]`, `["` + addrLower + `",{"blockHash":"` + hashUpper + `"}]`},
		{"eth_getTransactionCount", "address", `["` + addrLower + `","latest"]`, `["` + addrUpper + `","latest"]`},
		{"eth_getTransactionCount", "blockId", `["` + addrLower + `","0x1a"]`, `["` + addrLower + `","0x1A"]`},
		{"eth_getCode", "address", `["` + addrLower + `","latest"]`, `["` + addrUpper + `","latest"]`},
		{"eth_getCode", "blockId", `["` + addrLower + `","0x1a"]`, `["` + addrLower + `","0x1A"]`},
		// eth_call: txObject fields + blockId
		{"eth_call", "to", `[{"to":"` + addrLower + `"},"latest"]`, `[{"to":"` + addrUpper + `"},"latest"]`},
		{"eth_call", "from", `[{"from":"` + addrLower + `"},"latest"]`, `[{"from":"` + addrUpper + `"},"latest"]`},
		{"eth_call", "data", `[{"data":"0xabcd"},"latest"]`, `[{"data":"0xABCD"},"latest"]`},
		{"eth_call", "blockId", `[{},"0x1a"]`, `[{},"0x1A"]`},
		{"eth_call", "obj_blockHash", `[{},{"blockHash":"` + hashLower + `"}]`, `[{},{"blockHash":"` + hashUpper + `"}]`},
		// eth_getStorageAt: address + slot + blockId
		{"eth_getStorageAt", "address", `["` + addrLower + `","0x0","latest"]`, `["` + addrUpper + `","0x0","latest"]`},
		{"eth_getStorageAt", "slot", `["` + addrLower + `","0xab","latest"]`, `["` + addrLower + `","0xAB","latest"]`},
		{"eth_getStorageAt", "blockId", `["` + addrLower + `","0x0","0x1a"]`, `["` + addrLower + `","0x0","0x1A"]`},
		{"eth_getStorageAt", "obj_blockHash", `["` + addrLower + `","0x0",{"blockHash":"` + hashLower + `"}]`, `["` + addrLower + `","0x0",{"blockHash":"` + hashUpper + `"}]`},
		// eth_call: accessList nested address + storageKeys
		{"eth_call", "accessList_address", `[{"accessList":[{"address":"` + addrLower + `","storageKeys":[]}]},"latest"]`, `[{"accessList":[{"address":"` + addrUpper + `","storageKeys":[]}]},"latest"]`},
		{"eth_call", "accessList_storageKey", `[{"accessList":[{"address":"` + addrLower + `","storageKeys":["` + hashLower + `"]}]},"latest"]`, `[{"accessList":[{"address":"` + addrLower + `","storageKeys":["` + hashUpper + `"]}]},"latest"]`},
		// eth_call: stateOverrides address as object key
		{"eth_call", "override_address_key", `[{},"latest",{"` + addrLower + `":{"balance":"0x1"}}]`, `[{},"latest",{"` + addrUpper + `":{"balance":"0x1"}}]`},
		// eth_call: state slot as object key
		{"eth_call", "state_slot_key", `[{},"latest",{"` + addrLower + `":{"state":{"` + hashLower + `":"0x1"}}}]`, `[{},"latest",{"` + addrLower + `":{"state":{"` + hashUpper + `":"0x1"}}}]`},
	}
	for _, tt := range tests {
		t.Run(tt.method+"/"+tt.name, func(t *testing.T) {
			keyA := BuildKey(1, tt.method, mustParseParams(t, tt.a))
			keyB := BuildKey(1, tt.method, mustParseParams(t, tt.b))
			if keyA != keyB {
				t.Fatalf("case normalization failed:\n  %s → %s\n  %s → %s", tt.a, keyA, tt.b, keyB)
			}
		})
	}
}

func TestBuildKey_WhitespaceNormalization(t *testing.T) {
	compact := `[{"to":"0x0000000000000000000000000000000000000001","data":"0x0001"},"latest"]`
	spaced := `[ { "to" : "0x0000000000000000000000000000000000000001" , "data" : "0x0001" } , "latest" ]`

	keyA := BuildKey(1, "eth_call", mustParseParams(t, compact))
	keyB := BuildKey(1, "eth_call", mustParseParams(t, spaced))
	if keyA != keyB {
		t.Fatalf("whitespace normalization failed:\n  %s → %s\n  %s → %s", compact, keyA, spaced, keyB)
	}
}
