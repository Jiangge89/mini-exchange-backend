package com.miniexchange.txservice.domain;

import jakarta.persistence.*;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;

@Entity
@Table(name = "withdrawals")
public class Withdrawal {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "idempotency_id", unique = true, nullable = false)
    private String idempotencyId;

    @Column(name = "to_address", nullable = false)
    private String toAddress;

    @Column(nullable = false, precision = 78, scale = 0)
    private BigDecimal amount;

    private Integer nonce;

    @Column(name = "tx_hash")
    private String txHash;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TxStatus status = TxStatus.PENDING;

    @Column(name = "gas_price")
    private BigInteger gasPrice;

    @Column(name = "block_number")
    private Long blockNumber;

    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    @PrePersist
    void onCreate() {
        createdAt = Instant.now();
        updatedAt = Instant.now();
    }

    @PreUpdate
    void onUpdate() {
        updatedAt = Instant.now();
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public String getIdempotencyId() { return idempotencyId; }
    public void setIdempotencyId(String idempotencyId) { this.idempotencyId = idempotencyId; }

    public String getToAddress() { return toAddress; }
    public void setToAddress(String toAddress) { this.toAddress = toAddress; }

    public BigDecimal getAmount() { return amount; }
    public void setAmount(BigDecimal amount) { this.amount = amount; }

    public Integer getNonce() { return nonce; }
    public void setNonce(Integer nonce) { this.nonce = nonce; }

    public String getTxHash() { return txHash; }
    public void setTxHash(String txHash) { this.txHash = txHash; }

    public TxStatus getStatus() { return status; }
    public void setStatus(TxStatus status) { this.status = status; }

    public BigInteger getGasPrice() { return gasPrice; }
    public void setGasPrice(BigInteger gasPrice) { this.gasPrice = gasPrice; }

    public Long getBlockNumber() { return blockNumber; }
    public void setBlockNumber(Long blockNumber) { this.blockNumber = blockNumber; }

    public Instant getCreatedAt() { return createdAt; }
    public Instant getUpdatedAt() { return updatedAt; }
}
