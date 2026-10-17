package store

import (
	"testing"
)

func TestSaveBlock_WritesTransfersAndCheckpoint(t *testing.T) {
	t.Skip("TODO: implement")
	// 保存一个包含 2 条 Transfer 的区块
	// 期望数据库中有 2 条转账记录
	// 期望断点更新为新区块
	// 期望发送方余额减少，接收方余额增加
}

func TestSaveBlock_Idempotent(t *testing.T) {
	t.Skip("TODO: implement")
	// 同一个区块（相同 checkpoint + transfers）保存两次
	// 期望转账记录不重复（tx_hash + log_index 唯一约束）
	// 期望余额正确（不重复计算）
}

func TestSaveBlock_TransactionAtomicity(t *testing.T) {
	t.Skip("TODO: implement")
	// 模拟在写入过程中发生错误（如余额更新失败）
	// 期望整个事务回滚，转账记录也不写入
	// 期望断点不变
}

func TestSaveBlock_MintTransfer(t *testing.T) {
	t.Skip("TODO: implement")
	// from == address(0) 的转账
	// 期望只更新 to 的余额，不创建 address(0) 的余额记录
}

func TestRollbackToBlock_DeletesLaterBlocks(t *testing.T) {
	t.Skip("TODO: implement")
	// 数据库有区块 100-110 的数据
	// 回滚到区块 105
	// 期望 106-110 的转账记录被删除
	// 期望受影响地址的余额被重新计算
	// 期望断点回退到 105
}

func TestRollbackToBlock_RecalculatesBalances(t *testing.T) {
	t.Skip("TODO: implement")
	// 地址 A 在区块 100 收到 100 token，在区块 110 收到 50 token
	// 回滚到区块 105
	// 期望地址 A 的余额为 100（而非 150）
}

func TestGetTransfers_ByAddress(t *testing.T) {
	t.Skip("TODO: implement")
	// 数据库有多个地址的转账记录
	// 查询特定地址
	// 期望只返回该地址的发送和接收记录
}

func TestGetBalance_NoRecords(t *testing.T) {
	t.Skip("TODO: implement")
	// 查询没有任何记录的地址
	// 期望返回 0 而非错误
}
