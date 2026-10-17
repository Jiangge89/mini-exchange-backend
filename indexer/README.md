# indexer

链上 ERC-20 Transfer 事件索引器。

## 功能

- 扫描 ERC-20 合约的 Transfer 事件并写入 Postgres
- 断点续扫（重启后从上次位置继续）
- reorg 检测（校验 parentHash）与回滚
- 幂等写入（tx_hash + log_index 唯一约束）
- 余额计算与更新
- 每条记录标注确认级别（latest/safe/finalized）
- REST API 查询转账记录和余额

## 构建

```bash
go build ./cmd/indexer
```

## 运行

```bash
# 先执行数据库迁移
psql $DATABASE_URL < migrations/001_init.sql

# 启动
export TOKEN_CONTRACT_ADDRESS=0x...
export GETH_URL=http://localhost:8080
./indexer
```

## API

- `GET /api/transfers/{address}?limit=20&offset=0` — 查询地址的转账记录
- `GET /api/balance/{address}` — 查询地址的代币余额
- `GET /health` — 健康检查
- `GET /metrics` — Prometheus 指标

## 配置

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| INDEXER_PORT | 8081 | API 监听端口 |
| GETH_URL | http://localhost:8080 | 通过 gateway 访问链 |
| TOKEN_CONTRACT_ADDRESS | - | ERC-20 合约地址（必填） |
| POSTGRES_HOST | localhost | 数据库地址 |
| POSTGRES_PORT | 5432 | 数据库端口 |
| POSTGRES_USER | miniexchange | 数据库用户 |
| POSTGRES_PASSWORD | miniexchange | 数据库密码 |
| POSTGRES_DB | miniexchange | 数据库名 |
