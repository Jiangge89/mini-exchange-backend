package com.miniexchange.txservice.nonce;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class NonceManagerTest {

    @Test
    @Disabled("TODO: implement")
    void allocate_firstTime_queriesChainAndReturns() {
        // 首次为某地址分配 nonce
        // 期望从链上查询 eth_getTransactionCount 获取初始值
        // 期望返回该值并在数据库中记录 next_nonce = 初始值 + 1
    }

    @Test
    @Disabled("TODO: implement")
    void allocate_subsequentCalls_incrementsSequentially() {
        // 连续调用三次 allocate
        // 期望返回连续递增的 nonce（如 5, 6, 7）
    }

    @Test
    @Disabled("TODO: implement")
    void allocate_concurrentCalls_noCollision() {
        // 10 个线程同时调用 allocate
        // 期望返回 10 个不同的 nonce
        // 数据库行锁保证无竞态
    }

    @Test
    @Disabled("TODO: implement")
    void allocate_afterRestart_syncsWithChain() {
        // 数据库中 next_nonce = 5，但链上 transactionCount = 8
        // （重启期间有交易被确认）
        // 期望返回 max(5, 8) = 8，并更新数据库
    }

    @Test
    @Disabled("TODO: implement")
    void allocate_databaseLockTimeout_throwsException() {
        // 模拟行锁等待超时
        // 期望抛出异常而非无限等待
    }
}
