# gateway

区块链 JSON-RPC 网关。

## 功能

- 多上游节点管理（本机 geth、Alchemy、公共 RPC）
- 健康检查与自动故障切换
- 按 JSON-RPC 方法路由
- 感知最终确认的缓存（finalized 长缓存、latest 短缓存、写方法不缓存）
- 按 API key 限流
- 批量 JSON-RPC 请求支持
- Prometheus 指标

## 构建

```bash
go build ./cmd/gateway
```

## 运行

```bash
export GETH_URL=http://127.0.0.1:8545
export REDIS_URL=redis://localhost:6379
./gateway
```

## 测试

```bash
go test ./...
```

## 配置

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| GATEWAY_PORT | 8080 | 监听端口 |
| GETH_URL | - | 本机 geth 地址 |
| ALCHEMY_HOODI_URL | - | Alchemy Hoodi RPC |
| PUBLIC_RPC_URL | - | 公共 RPC 地址 |
| REDIS_URL | redis://localhost:6379 | Redis 地址 |

## 指标

- `gateway_requests_total{method}` — 请求总数（按方法）
- `gateway_request_duration_seconds{method}` — 请求延迟
- `gateway_upstream_errors_total` — 上游错误数
- `gateway_cache_hits_total` / `gateway_cache_misses_total` — 缓存命中/未命中
- `gateway_upstream_health_score{upstream}` — 上游健康分数
