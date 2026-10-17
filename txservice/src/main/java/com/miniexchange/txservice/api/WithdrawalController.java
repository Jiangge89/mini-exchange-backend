package com.miniexchange.txservice.api;

import com.miniexchange.txservice.domain.Withdrawal;
import com.miniexchange.txservice.domain.WithdrawalRepository;
import com.miniexchange.txservice.nonce.NonceManager;
import com.miniexchange.txservice.signer.TransactionSigner;
import com.miniexchange.txservice.broadcaster.TransactionBroadcaster;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.math.BigDecimal;

@RestController
@RequestMapping("/api/withdrawals")
public class WithdrawalController {

    private final WithdrawalRepository repository;
    private final NonceManager nonceManager;
    private final TransactionSigner signer;
    private final TransactionBroadcaster broadcaster;

    public WithdrawalController(
            WithdrawalRepository repository,
            NonceManager nonceManager,
            TransactionSigner signer,
            TransactionBroadcaster broadcaster) {
        this.repository = repository;
        this.nonceManager = nonceManager;
        this.signer = signer;
        this.broadcaster = broadcaster;
    }

    public record CreateWithdrawalRequest(
        String idempotencyId,
        String toAddress,
        BigDecimal amount
    ) {}

    public record WithdrawalResponse(
        Long id,
        String idempotencyId,
        String toAddress,
        String amount,
        String status,
        String txHash
    ) {}

    @PostMapping
    public ResponseEntity<WithdrawalResponse> createWithdrawal(@RequestBody CreateWithdrawalRequest request) {
        var existing = repository.findByIdempotencyId(request.idempotencyId());
        if (existing.isPresent()) {
            return ResponseEntity.ok(toResponse(existing.get()));
        }

        var withdrawal = new Withdrawal();
        withdrawal.setIdempotencyId(request.idempotencyId());
        withdrawal.setToAddress(request.toAddress());
        withdrawal.setAmount(request.amount());
        repository.save(withdrawal);

        return ResponseEntity.status(201).body(toResponse(withdrawal));
    }

    @GetMapping("/{id}")
    public ResponseEntity<WithdrawalResponse> getWithdrawal(@PathVariable Long id) {
        return repository.findById(id)
                .map(w -> ResponseEntity.ok(toResponse(w)))
                .orElse(ResponseEntity.notFound().build());
    }

    private static WithdrawalResponse toResponse(Withdrawal w) {
        return new WithdrawalResponse(
                w.getId(),
                w.getIdempotencyId(),
                w.getToAddress(),
                w.getAmount().toPlainString(),
                w.getStatus().name(),
                w.getTxHash()
        );
    }
}
