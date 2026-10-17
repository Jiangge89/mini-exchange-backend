package upstream

import (
	"context"
	"errors"
	"time"
)

// HealthChecker 定期检查所有上游节点的健康状态，并为每个节点计算健康分数。
//
// 健康分数计算因素：
// 1. 区块高度：与最高节点相比落后多少个区块。落后越多分数越低。
//    - 落后 0-1 个区块：不扣分
//    - 落后 2-5 个区块：轻微扣分
//    - 落后 >5 个区块：严重扣分
// 2. 错误率：最近 N 次请求中的失败比例。
//    - <1%：不扣分
//    - 1-10%：轻微扣分
//    - >10%：严重扣分
// 3. 延迟：平均响应时间。
//    - <100ms：不扣分
//    - 100ms-1s：轻微扣分
//    - >1s：严重扣分
//
// 边界情况：
// - 新加入的节点没有历史数据：给默认中等分数
// - 节点从不可用恢复：需要一段时间重建信任分数，避免抖动
// - 所有节点都不健康：仍然保留一个最佳的，不要全部标记为不可用
type HealthChecker interface {
	// Score 返回指定上游节点的健康分数。
	//
	// 返回值：0.0 到 1.0 之间的浮点数
	// - 1.0：完全健康
	// - 0.0：完全不可用
	//
	// 实现提示：
	// - 综合区块高度差、错误率、延迟三个维度
	// - 可以使用加权平均或乘法模型
	// - 考虑时间衰减：最近的数据权重更高
	Score(u *Upstream) float64

	// Start 启动后台健康检查循环。
	//
	// 行为：
	// - 每隔固定间隔（如 5 秒）轮询每个上游的 eth_blockNumber
	// - 更新各节点的 blockHeight
	// - 检测到节点无响应时记录错误
	// - ctx 取消时停止
	Start(ctx context.Context)
}

type healthChecker struct {
	pool     *Pool
	interval time.Duration
}

func NewHealthChecker(pool *Pool) HealthChecker {
	return &healthChecker{
		pool:     pool,
		interval: 5 * time.Second,
	}
}

func (h *healthChecker) Score(_ *Upstream) float64 {
	panic(errors.New("not implemented"))
}

func (h *healthChecker) Start(_ context.Context) {
	panic(errors.New("not implemented"))
}
