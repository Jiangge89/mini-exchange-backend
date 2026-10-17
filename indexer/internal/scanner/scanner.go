package scanner

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/chain"
	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/store"
)

// Scanner 是链上事件索引器的扫描主循环。
//
// 职责：
// 1. 从数据库读取上次扫描到的断点（区块号 + 区块哈希）
// 2. 获取下一个区块的区块头
// 3. 校验 parentHash 是否匹配断点的区块哈希（reorg 检测）
// 4. 如果不匹配，说明发生了链重组，需要回滚到分叉点
// 5. 如果匹配，拉取该区块的 Transfer 日志
// 6. 在一个数据库事务中写入转账记录、更新余额、更新断点
// 7. 重复以上步骤
type Scanner struct {
	client        *chain.Client
	store         store.Store
	tokenContract common.Address
}

func New(client *chain.Client, store store.Store, tokenContract common.Address) *Scanner {
	return &Scanner{
		client:        client,
		store:         store,
		tokenContract: tokenContract,
	}
}

// Run 启动扫描主循环。
//
// 行为：
// - 持续循环扫描新区块，直到 ctx 被取消
// - 每次循环：
//   1. 读取断点（上次扫描到的区块号和哈希）
//   2. 获取下一个区块头（断点区块号 + 1）
//   3. 校验 parentHash == 断点哈希
//      - 不匹配 → 调用 handleReorg 回滚
//      - 匹配 → 调用 processBlock 处理区块
//   4. 如果当前追上链头，等待一段时间再重试
//
// 参数：
// - ctx: 取消时优雅停止循环
//
// 边界情况：
// - 首次启动，断点不存在：从合约部署区块或区块 0 开始
// - 链头暂时没有新区块：等待后重试，不要消耗 CPU
// - 中途 ctx 取消：完成当前区块的处理后退出，不要写一半
// - 网关返回错误：指数退避重试
// - 重启后断点存在但区块已经过了很多：逐块追赶，不要跳过
func (s *Scanner) Run(ctx context.Context) {
	panic(errors.New("not implemented"))
}

// handleReorg 检测链重组并回滚到分叉点。
//
// 参数：
// - ctx: 上下文
// - currentHeight: 当前断点的区块号
// - expectedParentHash: 断点记录的区块哈希
//
// 行为：
// 1. 从 currentHeight 开始向前查找，逐块比较 parentHash
// 2. 找到分叉点（链上和数据库一致的最近祖先）
// 3. 删除分叉点之后的所有转账记录
// 4. 重新计算受影响地址的余额
// 5. 将断点回退到分叉点
//
// 边界情况：
// - 深度 reorg（>10 个区块）：仍然处理，但记录告警日志
// - 回滚中途出错：事务回滚，下次重试
// - reorg 之后又 reorg：正常处理，因为每次都重新校验 parentHash
func (s *Scanner) handleReorg(ctx context.Context, currentHeight uint64, expectedParentHash common.Hash) error {
	return errors.New("not implemented")
}

// processBlock 处理单个区块：拉取日志、写入数据库。
//
// 参数：
// - ctx: 上下文
// - blockNumber: 要处理的区块号
//
// 行为：
// 1. 用 eth_getLogs 拉取该区块中目标合约的 Transfer 事件
// 2. 在一个数据库事务中：
//    a. 写入每条转账记录（tx_hash + log_index 作为唯一键，实现幂等）
//    b. 更新发送方和接收方的余额
//    c. 更新断点为当前区块号和哈希
//    d. 标注每条记录的确认级别（latest/safe/finalized）
// 3. 提交事务
//
// 边界情况：
// - 区块中没有 Transfer 事件：仍然更新断点
// - 重复处理同一个区块（幂等）：INSERT ON CONFLICT DO NOTHING
// - 数据库事务失败：整个区块回滚，下次重试
// - Transfer 事件中 from == address(0)：这是 mint 操作，只增加 to 的余额
// - Transfer 事件中 to == address(0)：这是 burn 操作，只减少 from 的余额
func (s *Scanner) processBlock(ctx context.Context, blockNumber uint64) error {
	return errors.New("not implemented")
}
