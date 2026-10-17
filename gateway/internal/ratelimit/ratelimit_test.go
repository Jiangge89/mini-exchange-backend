package ratelimit

import (
	"testing"
)

func TestAllow_UnderLimit_Allowed(t *testing.T) {
	t.Skip("TODO: implement")
	// API key 还在配额内
	// 期望 allowed == true, remaining > 0
}

func TestAllow_AtLimit_Rejected(t *testing.T) {
	t.Skip("TODO: implement")
	// API key 已达到每秒限制
	// 期望 allowed == false, remaining == 0
}

func TestAllow_WindowExpiry_Resets(t *testing.T) {
	t.Skip("TODO: implement")
	// API key 达到限制后等待窗口过期
	// 期望新窗口内 allowed == true
}

func TestAllow_DifferentKeys_Independent(t *testing.T) {
	t.Skip("TODO: implement")
	// key-A 达到限制，key-B 不应受影响
}

func TestAllow_RedisDown_FailOpen(t *testing.T) {
	t.Skip("TODO: implement")
	// Redis 不可用时
	// 期望 allowed == true（fail-open 策略）
}

func TestAllow_EmptyKey_Handled(t *testing.T) {
	t.Skip("TODO: implement")
	// 空的 API key
	// 期望有合理的行为（全局限流或放行）
}

func TestAllow_ConcurrentRequests_Atomic(t *testing.T) {
	t.Skip("TODO: implement")
	// 100 个并发请求使用同一个 key
	// 期望总放行数不超过限额（验证原子性）
}
