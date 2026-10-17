package com.miniexchange.txservice.domain;

public enum TxStatus {
    /** 提现请求已创建，等待处理 */
    PENDING,
    /** 交易已广播到链上 */
    BROADCAST,
    /** 交易已被区块包含 */
    INCLUDED,
    /** 交易已被最终确认（不会因 reorg 回退） */
    FINALIZED,
    /** 交易失败（执行 revert 或 gas 不足等） */
    FAILED,
    /** 交易已被替换（同 nonce 加价重发） */
    REPLACED
}
