package com.miniexchange.txservice.nonce;

import org.springframework.stereotype.Component;

/**
 * NonceManager 负责为交易分配 nonce，使用数据库行锁保证多实例下的唯一性。
 *
 * <p>设计要点：
 * <ul>
 *   <li>每个发送地址维护一个 nonce 计数器（存储在数据库中）</li>
 *   <li>使用 SELECT ... FOR UPDATE 行锁保证并发安全</li>
 *   <li>分配后 nonce 自增，即使交易最终失败也不回退 nonce</li>
 * </ul>
 *
 * <p>边界情况：
 * <ul>
 *   <li>首次使用某地址：需要从链上查询当前 nonce 作为起始值</li>
 *   <li>nonce 间隙（gap）：如果某个 nonce 的交易失败，后续 nonce 的交易会卡住。
 *       需要检测并用空交易填补间隙，或回收该 nonce</li>
 *   <li>重启后 nonce 同步：数据库中的 nonce 可能落后于链上（有交易在重启期间被确认），
 *       需要从链上获取最新 nonce 并取两者的较大值</li>
 *   <li>多实例竞争：两个实例同时请求 nonce，行锁保证只有一个能成功，另一个等待</li>
 *   <li>数据库连接超时：行锁等待时间过长时应超时失败而非无限等待</li>
 * </ul>
 */
@Component
public class NonceManager {

    /**
     * 为指定发送地址分配一个 nonce。
     *
     * <p>实现步骤：
     * <ol>
     *   <li>开启数据库事务</li>
     *   <li>SELECT next_nonce FROM nonce_counter WHERE address = ? FOR UPDATE（行锁）</li>
     *   <li>如果不存在，从链上查询 eth_getTransactionCount 初始化</li>
     *   <li>分配当前 nonce</li>
     *   <li>UPDATE nonce_counter SET next_nonce = next_nonce + 1</li>
     *   <li>提交事务（释放锁）</li>
     * </ol>
     *
     * @param senderAddress 发送交易的以太坊地址
     * @return 分配的 nonce 值
     * @throws RuntimeException 如果数据库锁等待超时或链上查询失败
     */
    public int allocate(String senderAddress) {
        throw new UnsupportedOperationException("not implemented");
    }
}
