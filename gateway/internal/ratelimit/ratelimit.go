package ratelimit

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// RateLimiter 按 API key 对请求进行限流。
//
// 设计要点：
// - 使用 Redis 实现，支持多实例网关共享限流状态
// - 每个 API key 有独立的配额
//
// 推荐算法（任选其一）：
// 1. 令牌桶（Token Bucket）：允许一定的突发流量
// 2. 滑动窗口计数器（Sliding Window Counter）：更精确，实现稍复杂
// 3. 固定窗口计数器（Fixed Window Counter）：最简单，但在窗口边界有突发问题
//
// 边界情况：
// - Redis 不可用时的降级策略：应该放行还是拒绝？建议放行（fail-open）
// - API key 为空时：可以用一个全局共享的限流桶，或直接放行
// - 时钟漂移：如果用时间窗口，需要考虑多实例间的时钟差异
// - 原子性：检查和扣减必须是原子操作（用 Redis Lua 脚本或 MULTI/EXEC）
type RateLimiter interface {
	// Allow 检查给定 API key 的请求是否被允许。
	//
	// 参数：
	// - ctx: 上下文
	// - apiKey: 客户端的 API key
	//
	// 返回值：
	// - allowed: 是否允许此请求
	// - remaining: 当前窗口内剩余可用请求数
	// - err: Redis 错误等。注意：错误时建议 fail-open（返回 allowed=true）
	//
	// 实现提示：
	// - 默认限制：每个 key 每秒 50 请求（可配置）
	// - 使用 Redis 的 INCR + EXPIRE 或 Lua 脚本实现原子操作
	// - key 格式建议：ratelimit:{apiKey}:{window}
	Allow(ctx context.Context, apiKey string) (allowed bool, remaining int, err error)
}

type rateLimiter struct {
	client *redis.Client
}

func New(redisURL string) RateLimiter {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		opt = &redis.Options{Addr: "localhost:6379"}
	}
	return &rateLimiter{client: redis.NewClient(opt)}
}

func (r *rateLimiter) Allow(_ context.Context, _ string) (bool, int, error) {
	return true, 0, errors.New("not implemented")
}
