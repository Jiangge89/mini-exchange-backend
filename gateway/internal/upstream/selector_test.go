package upstream

import (
	"testing"
)

func TestSelector_Pick_ReturnsHealthiestUpstream(t *testing.T) {
	t.Skip("TODO: implement")
	// 准备三个上游，健康分数分别为 0.9、0.5、0.2
	// 期望 Pick() 返回分数 0.9 的上游
}

func TestSelector_Pick_AllUnhealthy_ReturnsBestEffort(t *testing.T) {
	t.Skip("TODO: implement")
	// 所有上游分数都很低（如 0.1）
	// 期望 Pick() 仍然返回分数最高的那个，而不是报错
}

func TestSelector_Pick_EmptyPool_ReturnsError(t *testing.T) {
	t.Skip("TODO: implement")
	// 上游池为空
	// 期望返回 ErrNoHealthyUpstream
}

func TestSelector_Forward_Success(t *testing.T) {
	t.Skip("TODO: implement")
	// 模拟一个健康的上游
	// 发送请求，期望收到正确的响应
	// 验证 RecordSuccess 被调用
}

func TestSelector_Forward_FailoverToNext(t *testing.T) {
	t.Skip("TODO: implement")
	// 第一个上游返回错误，第二个上游正常
	// 期望自动切换到第二个上游
	// 验证第一个上游的 RecordError 被调用
}

func TestSelector_Forward_AllFail(t *testing.T) {
	t.Skip("TODO: implement")
	// 所有上游都返回错误
	// 期望返回错误，不要无限重试
}

func TestSelector_ConcurrentPick(t *testing.T) {
	t.Skip("TODO: implement")
	// 100 个 goroutine 并发调用 Pick()
	// 期望没有 race condition（用 -race 检测）
}
