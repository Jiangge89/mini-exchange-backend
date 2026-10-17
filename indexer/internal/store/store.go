package store

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type ConfirmationLevel string

const (
	LevelLatest    ConfirmationLevel = "latest"
	LevelSafe      ConfirmationLevel = "safe"
	LevelFinalized ConfirmationLevel = "finalized"
)

type Transfer struct {
	ID                int64
	BlockNumber       uint64
	BlockHash         common.Hash
	TxHash            common.Hash
	LogIndex          uint
	From              common.Address
	To                common.Address
	Value             *big.Int
	ConfirmationLevel ConfirmationLevel
	CreatedAt         time.Time
}

type Balance struct {
	Address common.Address
	Amount  *big.Int
}

type Checkpoint struct {
	BlockNumber uint64
	BlockHash   common.Hash
}

// Store 定义索引器的数据存储接口。
type Store interface {
	// GetCheckpoint 获取上次扫描到的断点。
	// 如果没有断点记录（首次启动），返回 nil, nil。
	GetCheckpoint(ctx context.Context) (*Checkpoint, error)

	// SaveBlock 在一个数据库事务中保存单个区块的所有数据。
	//
	// 这是核心的幂等写入操作，需要你实现。
	//
	// 参数：
	// - ctx: 上下文
	// - checkpoint: 新的断点（当前区块号和哈希）
	// - transfers: 该区块中的所有 Transfer 事件
	//
	// 行为（在单个事务中）：
	// 1. 使用 INSERT ON CONFLICT (tx_hash, log_index) DO NOTHING 写入转账记录
	// 2. 更新每个涉及地址的余额（发送方减少，接收方增加）
	//    - 使用 INSERT ON CONFLICT (address) DO UPDATE SET amount = amount +/- value
	//    - mint (from=0x0): 只增加 to 的余额
	//    - burn (to=0x0): 只减少 from 的余额
	// 3. 更新断点
	// 4. 提交事务
	//
	// 幂等性保证：
	// - 重复调用相同的 checkpoint + transfers 不会产生重复数据
	// - 余额计算基于实际写入的记录数，而非参数中的记录数
	//
	// 失败场景：
	// - 事务中任何一步失败，整个事务回滚
	// - 返回错误让调用方重试
	SaveBlock(ctx context.Context, checkpoint Checkpoint, transfers []Transfer) error

	// RollbackToBlock 回滚到指定区块（用于 reorg 处理）。
	//
	// 行为：
	// 1. 删除 blockNumber 之后的所有转账记录
	// 2. 重新计算受影响地址的余额
	// 3. 将断点回退到 blockNumber
	//
	// 参数：
	// - blockNumber: 回滚到的目标区块号（该区块保留）
	//
	// 边界情况：
	// - blockNumber 大于当前最大区块号：什么都不做
	// - blockNumber 为 0：清空所有数据
	RollbackToBlock(ctx context.Context, blockNumber uint64) error

	// GetTransfers 查询转账记录（供 REST API 使用）。
	GetTransfers(ctx context.Context, address common.Address, limit, offset int) ([]Transfer, error)

	// GetBalance 查询某个地址的代币余额。
	GetBalance(ctx context.Context, address common.Address) (*big.Int, error)

	// Close 关闭数据库连接。
	Close()
}
