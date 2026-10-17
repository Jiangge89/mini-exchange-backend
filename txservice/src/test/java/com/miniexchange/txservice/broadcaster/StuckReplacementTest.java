package com.miniexchange.txservice.broadcaster;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

class StuckReplacementTest {

    @Test
    @Disabled("TODO: implement")
    void replaceStuck_transactionStuckTooLong_replacesWithHigherGas() {
        // 交易广播后超过 30 个区块仍未被包含
        // 期望用相同 nonce + 1.1x gas price 重新发送
        // 期望新的 txHash 记录到数据库
    }

    @Test
    @Disabled("TODO: implement")
    void replaceStuck_gasPrice_atLeast10PercentHigher() {
        // 原交易 gasPrice = 100 gwei
        // 期望替换交易 gasPrice >= 110 gwei
    }

    @Test
    @Disabled("TODO: implement")
    void replaceStuck_exceedsMaxGasPrice_marksFailed() {
        // 多次替换后 gas price 超过上限
        // 期望标记为 FAILED 而非继续加价
    }

    @Test
    @Disabled("TODO: implement")
    void replaceStuck_nonceTooLow_investigates() {
        // 替换时节点返回 "nonce too low"
        // 说明原交易已经被确认，或另一笔交易用了这个 nonce
        // 期望正确处理（查询回执确认状态）
    }

    @Test
    @Disabled("TODO: implement")
    void replaceStuck_replacementUnderpriced_retriesWithHigherGas() {
        // 节点返回 "replacement transaction underpriced"
        // 期望提高加价幅度后重试
    }

    @Test
    @Disabled("TODO: implement")
    void broadcastPending_idempotent_skipsAlreadyBroadcasted() {
        // 同一个 PENDING 提现不应被重复广播
        // 期望只处理一次
    }
}
