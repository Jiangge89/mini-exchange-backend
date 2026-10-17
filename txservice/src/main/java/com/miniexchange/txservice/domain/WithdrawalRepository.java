package com.miniexchange.txservice.domain;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface WithdrawalRepository extends JpaRepository<Withdrawal, Long> {

    Optional<Withdrawal> findByIdempotencyId(String idempotencyId);

    List<Withdrawal> findByStatus(TxStatus status);

    @Query("SELECT w FROM Withdrawal w WHERE w.status = 'BROADCAST' AND w.updatedAt < :cutoff")
    List<Withdrawal> findStuckBroadcasts(java.time.Instant cutoff);

    List<Withdrawal> findByStatusAndBlockNumberGreaterThan(TxStatus status, Long blockNumber);
}
