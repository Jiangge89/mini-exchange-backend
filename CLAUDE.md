# CLAUDE.md — mini-exchange-backend 协作约定

## 项目背景

迷你交易所链上后端，运行在以太坊 Hoodi 测试网 (chainId 560048)，由三个服务组成：
- **gateway** (Go)：区块链 JSON-RPC 网关，多上游故障切换，感知最终确认的缓存，按 API key 限流
- **indexer** (Go)：链上 ERC-20 Transfer 事件索引器，断点续扫，reorg 检测与回滚，幂等写入
- **txservice** (Java/Spring Boot)：交易发送服务，模拟交易所提现，nonce 管理，卡单替换，状态机追踪

目标：作为开源作品集项目，展示对区块链后端核心难题（reorg、nonce、幂等、缓存失效）的理解。

## 分工约定

### 核心逻辑由我（人类）编写

标有 `TODO: implement` 的核心函数是留给我亲手实现的。除非我明确要求，**不要实现这些函数**。它们包括：

**gateway:**
- 上游选择与故障切换 (`selector.go`)
- 健康评分 (`health.go`)
- 缓存策略 (`cache.go`)
- 限流算法 (`ratelimit.go`)

**indexer:**
- 扫描主循环 (`scanner.go`)
- reorg 检测与回滚
- 单个区块的幂等写入 (`store.go` 中的事务写入)

**txservice:**
- `NonceManager.allocate` — 数据库行锁分配 nonce
- 卡单替换策略 (`TransactionBroadcaster`)
- 确认追踪与状态机流转 (`ConfirmationTracker`)
- 提现请求幂等校验

### 策略决定也由我做

- 涉及策略的决定（身份识别、限流对象、失败时放行还是拒绝、缓存什么不缓存什么）也由我决定。遇到需要这类决定的地方，留 TODO 注释说明需要决定什么，不要自行选择默认行为。
- 不要吞掉错误；配置错误在启动时快速失败。

### 你可以放手写

- 样板代码、配置加载、HTTP 路由、数据库连接
- 测试骨架和辅助测试工具
- Dockerfile、CI、文档、Makefile
- 新的工具函数或辅助模块

### Review 约定

当我请你 review 我的实现时，重点检查：
- 并发安全性（锁、channel、goroutine 泄漏）
- 幂等性（重复请求、重复事件）
- reorg 处理（回滚深度、parentHash 校验）
- nonce 管理（竞态、间隙、替换）
- 缓存失效（finalized vs latest、reorg 后失效）

指出问题并解释原因，但不要直接替我重写核心逻辑。

## 安全约定

- 仅用于 Hoodi 测试网，不用于主网
- 所有密钥、私钥只从环境变量读取
- `.env` 文件不提交（已在 `.gitignore` 中排除）
- 日志中永远不打印私钥

## 构建与测试命令

```bash
# 全局
make build          # 构建所有服务
make test           # 运行所有测试
make lint           # 代码检查
make up / down      # 启停 Docker Compose 基础设施

# gateway (Go)
cd gateway
go build ./...
go test ./...
go vet ./...

# indexer (Go)
cd indexer
go build ./...
go test ./...
go vet ./...

# txservice (Java)
cd txservice
./gradlew build
./gradlew test
./gradlew bootRun

# contracts (Solidity)
cd contracts
forge build
forge test
```
