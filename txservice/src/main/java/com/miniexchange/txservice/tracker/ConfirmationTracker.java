package com.miniexchange.txservice.tracker;

import org.springframework.stereotype.Component;

/**
 * ConfirmationTracker 追踪已广播交易的确认状态，驱动状态机流转。
 *
 * <p>状态机：
 * <pre>
 * PENDING → BROADCAST → INCLUDED → FINALIZED
 *                  ↘                 ↗
 *                   → REPLACED →  (同 nonce 新交易走 BROADCAST → ...)
 *                  ↘
 *                   → FAILED
 * </pre>
 *
 * <p>确认追踪流程：
 * <ol>
 *   <li>定时扫描所有 BROADCAST 状态的提现</li>
 *   <li>用 eth_getTransactionReceipt 查询交易回执</li>
 *   <li>有回执 → 检查 status：
 *     <ul>
 *       <li>status == 1（成功）→ 更新为 INCLUDED，记录 blockNumber</li>
 *       <li>status == 0（失败）→ 更新为 FAILED</li>
 *     </ul>
 *   </li>
 *   <li>无回执 → 交易仍在 pending，不做处理</li>
 * </ol>
 *
 * <p>Finalized 确认：
 * <ol>
 *   <li>扫描所有 INCLUDED 状态的提现</li>
 *   <li>查询链上 finalized 区块号</li>
 *   <li>如果提现的 blockNumber <= finalized 区块号 → 更新为 FINALIZED</li>
 * </ol>
 *
 * <p>Reorg 回退（核心逻辑，需要你实现）：
 * <ul>
 *   <li>如果 INCLUDED 的交易所在区块被 reorg 了：
 *     <ul>
 *       <li>回执查询返回 null（交易不再在链上）</li>
 *       <li>需要将状态从 INCLUDED 回退到 BROADCAST</li>
 *       <li>交易仍然在 mempool 中，可能被重新包含</li>
 *     </ul>
 *   </li>
 *   <li>如果 reorg 深度很大，可能需要重新广播</li>
 *   <li>已经 FINALIZED 的交易不会被 reorg</li>
 * </ul>
 *
 * <p>边界情况：
 * <ul>
 *   <li>节点暂时不可用：跳过本轮追踪，下次重试</li>
 *   <li>交易被替换后原 txHash 永远没有回执：靠 replaceStuck 处理</li>
 *   <li>同一笔提现有多个 txHash（替换产生的）：只追踪最新的 txHash</li>
 *   <li>区块确认数回退（reorg）：INCLUDED → BROADCAST 的回退</li>
 * </ul>
 */
@Component
public class ConfirmationTracker {

    /**
     * 追踪 BROADCAST 状态交易的确认情况。
     *
     * <p>对每个 BROADCAST 的提现，查询 eth_getTransactionReceipt：
     * <ul>
     *   <li>有回执且成功 → INCLUDED</li>
     *   <li>有回执且失败 → FAILED</li>
     *   <li>无回执 → 保持 BROADCAST</li>
     * </ul>
     */
    public void trackBroadcasted() {
        throw new UnsupportedOperationException("not implemented");
    }

    /**
     * 检查 INCLUDED 的交易是否达到 finalized。
     *
     * <p>查询链上 finalized 区块号，将 blockNumber <= finalized 的
     * INCLUDED 交易更新为 FINALIZED。
     */
    public void checkFinalized() {
        throw new UnsupportedOperationException("not implemented");
    }

    /**
     * 检测 reorg 并回退受影响的交易状态。
     *
     * <p>对 INCLUDED 状态的交易重新查询回执，如果回执不存在，
     * 说明该交易所在的区块被 reorg 了，需要回退到 BROADCAST 状态。
     */
    public void detectReorg() {
        throw new UnsupportedOperationException("not implemented");
    }
}
