package com.miniexchange.txservice.api;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

class WithdrawalIdempotencyTest {

    @Test
    @Disabled("TODO: implement")
    void createWithdrawal_duplicateIdempotencyId_returnsExisting() {
        // 用同一个 idempotencyId 发送两次请求
        // 第一次返回 201 Created
        // 第二次返回 200 OK，内容与第一次相同
        // 数据库中只有一条记录
    }

    @Test
    @Disabled("TODO: implement")
    void createWithdrawal_differentIdempotencyId_createsBoth() {
        // 用不同的 idempotencyId 发送两次请求
        // 期望创建两条独立的提现记录
    }

    @Test
    @Disabled("TODO: implement")
    void createWithdrawal_concurrentDuplicates_onlyOneCreated() {
        // 10 个线程同时用相同 idempotencyId 发送请求
        // 期望只有一条记录被创建
        // 数据库唯一约束保证
    }

    @Test
    @Disabled("TODO: implement")
    void createWithdrawal_invalidAddress_returns400() {
        // 无效的以太坊地址
        // 期望返回 400 Bad Request
    }

    @Test
    @Disabled("TODO: implement")
    void getWithdrawal_exists_returns200() {
        // 查询存在的提现
        // 期望返回完整信息
    }

    @Test
    @Disabled("TODO: implement")
    void getWithdrawal_notExists_returns404() {
        // 查询不存在的 ID
        // 期望返回 404
    }
}
