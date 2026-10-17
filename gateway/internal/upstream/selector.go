package upstream

import (
	"errors"
)

// Selector 负责从多个上游节点中选择一个可用节点，并在失败时切换到下一个。
//
// 设计要点：
// - Pick 从候选上游中选出健康分数最高的节点
// - Forward 将请求发送到选中的节点，如果失败则依次尝试其他节点
// - 每次失败都应调用 Upstream.RecordError()，成功调用 RecordSuccess()
//
// 边界情况：
// - 所有节点都不可用：返回明确的 "no healthy upstream" 错误
// - 只剩一个节点时的降级：即使不健康也应尝试
// - 并发安全：多个请求同时调用 Pick/Forward
type Selector interface {
	// Pick 选择一个可用的上游节点。
	//
	// 返回值：
	// - *Upstream: 选中的上游节点
	// - error: 所有节点不可用时返回 ErrNoHealthyUpstream
	//
	// 实现提示：
	// - 根据 HealthChecker.Score() 排序，优先选分数最高的
	// - 如果所有节点分数为 0，仍应返回一个（降级模式）
	Pick() (*Upstream, error)

	// Forward 将原始 JSON-RPC 请求转发到上游节点，失败时自动切换。
	//
	// 参数：
	// - ctx: 请求上下文，用于取消和超时
	// - reqBody: 完整的 JSON-RPC 请求体（已序列化的 JSON bytes）
	//
	// 返回值：
	// - []byte: 上游的原始响应体
	// - error: 所有上游都失败时返回错误
	//
	// 实现提示：
	// - 先调用 Pick() 选择节点，发送 HTTP POST 请求
	// - 如果请求失败（网络错误或 HTTP 5xx），ReportFailure 后尝试下一个节点
	// - 成功时调用 ReportSuccess 并记录延迟
	// - 最多重试 len(upstreams) 次
	// - 注意：JSON-RPC 层面的错误（如 method not found）不算上游失败
	Forward(ctx interface{ Done() <-chan struct{} }, reqBody []byte) ([]byte, error)
}

var ErrNoHealthyUpstream = errors.New("no healthy upstream available")

type selector struct {
	pool          *Pool
	healthChecker HealthChecker
}

func NewSelector(pool *Pool) Selector {
	return &selector{
		pool:          pool,
		healthChecker: nil, // 将在 main.go 中注入
	}
}

func (s *selector) Pick() (*Upstream, error) {
	return nil, errors.New("not implemented")
}

func (s *selector) Forward(ctx interface{ Done() <-chan struct{} }, reqBody []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}
