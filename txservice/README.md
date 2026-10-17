# txservice

交易发送服务，模拟交易所提现。

## 功能

- 提现请求幂等（业务幂等 ID）
- 多实例下使用数据库行锁分配 nonce
- 手续费估算
- 卡单检测与同 nonce 加价替换
- 交易状态机：PENDING → BROADCAST → INCLUDED → FINALIZED / FAILED / REPLACED
- 确认追踪与 reorg 回退

## 构建

```bash
./gradlew build
```

## 运行

```bash
export SIGNER_PRIVATE_KEY=0x...
export GETH_URL=http://localhost:8080
./gradlew bootRun
```

## 测试

```bash
./gradlew test
```

## API

- `POST /api/withdrawals` — 创建提现请求
  ```json
  {
    "idempotencyId": "withdraw-001",
    "toAddress": "0x...",
    "amount": "1000000000000000000"
  }
  ```
- `GET /api/withdrawals/{id}` — 查询提现状态
- `GET /actuator/health` — 健康检查
- `GET /actuator/prometheus` — Prometheus 指标

## 配置

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| TXSERVICE_PORT | 8082 | 监听端口 |
| GETH_URL | http://localhost:8080 | 通过 gateway 访问链 |
| SIGNER_PRIVATE_KEY | - | 签名私钥（仅测试网！必填） |
| CHAIN_ID | 560048 | 链 ID（Hoodi） |
| POSTGRES_HOST | localhost | 数据库地址 |
| POSTGRES_PORT | 5432 | 数据库端口 |
| POSTGRES_USER | miniexchange | 数据库用户 |
| POSTGRES_PASSWORD | miniexchange | 数据库密码 |
| POSTGRES_DB | miniexchange | 数据库名 |

## 状态机

```
PENDING ──→ BROADCAST ──→ INCLUDED ──→ FINALIZED
                │                        ↑
                ├──→ REPLACED ──→ (新交易走同样的流程)
                │
                └──→ FAILED
```
