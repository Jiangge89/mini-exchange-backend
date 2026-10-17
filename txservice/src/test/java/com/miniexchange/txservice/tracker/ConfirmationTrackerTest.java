package com.miniexchange.txservice.tracker;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

class ConfirmationTrackerTest {

    @Test
    @Disabled("TODO: implement")
    void trackBroadcasted_receiptFound_success_updatesToIncluded() {
        // 交易有回执且 status == 1
        // 期望状态从 BROADCAST → INCLUDED
        // 期望记录 blockNumber
    }

    @Test
    @Disabled("TODO: implement")
    void trackBroadcasted_receiptFound_failed_updatesToFailed() {
        // 交易有回执且 status == 0（revert）
        // 期望状态从 BROADCAST → FAILED
    }

    @Test
    @Disabled("TODO: implement")
    void trackBroadcasted_noReceipt_staysBroadcasted() {
        // 交易没有回执（还在 mempool）
        // 期望状态保持 BROADCAST
    }

    @Test
    @Disabled("TODO: implement")
    void checkFinalized_blockFinalized_updatesToFinalized() {
        // 交易在区块 100，链上 finalized 区块 = 105
        // 期望状态从 INCLUDED → FINALIZED
    }

    @Test
    @Disabled("TODO: implement")
    void checkFinalized_blockNotYetFinalized_staysIncluded() {
        // 交易在区块 100，链上 finalized 区块 = 95
        // 期望状态保持 INCLUDED
    }

    @Test
    @Disabled("TODO: implement")
    void detectReorg_receiptDisappears_rollsBackToBroadcasted() {
        // INCLUDED 的交易重新查询回执返回 null
        // 说明区块被 reorg 了
        // 期望状态从 INCLUDED → BROADCAST
        // 期望清除 blockNumber
    }

    @Test
    @Disabled("TODO: implement")
    void detectReorg_finalizedTransaction_notAffected() {
        // FINALIZED 的交易不应被 reorg 检查影响
        // 期望状态保持 FINALIZED
    }

    @Test
    @Disabled("TODO: implement")
    void trackBroadcasted_nodeUnavailable_skipsGracefully() {
        // 节点返回连接错误
        // 期望不崩溃，日志记录错误，下次重试
    }
}
