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

func TestBuildKey_Deterministic(t *testing.T) {
	params := []json.RawMessage{json.RawMessage(`"latest"`), json.RawMessage(`true`)}
	key1 := BuildKey("eth_getBlockByNumber", params)
	key2 := BuildKey("eth_getBlockByNumber", params)
	if key1 != key2 {
		t.Errorf("BuildKey should be deterministic, got %s and %s", key1, key2)
	}
}

func TestBuildKey_DifferentMethodsDifferentKeys(t *testing.T) {
	params := []json.RawMessage{json.RawMessage(`"latest"`)}
	key1 := BuildKey("eth_getBlockByNumber", params)
	key2 := BuildKey("eth_getBalance", params)
	if key1 == key2 {
		t.Error("different methods should produce different keys")
	}
}
