package upstream

import (
	"testing"
)

func TestHealthChecker_Score_FullyHealthy(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游区块高度最新、无错误、低延迟
	// 期望 Score ≈ 1.0
}

func TestHealthChecker_Score_BlockHeightBehind(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游区块高度落后 3 个区块
	// 期望 Score 在 0.5-0.9 之间（部分扣分）
}

func TestHealthChecker_Score_HighErrorRate(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游最近 100 次请求有 20 次失败
	// 期望 Score 明显降低
}

func TestHealthChecker_Score_HighLatency(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游平均延迟 2s
	// 期望 Score 降低
}

func TestHealthChecker_Score_NewUpstream_DefaultScore(t *testing.T) {
	t.Skip("TODO: implement")
	// 新上游没有任何历史数据
	// 期望返回一个合理的默认分数（如 0.5）
}

func TestHealthChecker_Score_Recovery(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游之前错误率很高，最近一段时间恢复正常
	// 期望分数逐步回升，而不是立刻跳回 1.0（避免抖动）
}
