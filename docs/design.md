# Mini Exchange Backend — Design Decisions

This document explains the **why** behind each core design decision. It is written for reviewers and interviewers who want to understand the tradeoffs, not just read the code.

---

## Architecture Overview

Three services, one ERC-20 token contract, one testnet (Hoodi, chainId 560048).

```
                        ┌────────────┐
                        │  MiniToken │  ERC-20 on Hoodi
                        │  (Solidity)│
                        └──────┬─────┘
                               │
            ┌──────────────────┼──────────────────┐
            │                  │                  │
      ┌─────┴─────┐    ┌──────┴──────┐    ┌──────┴──────┐
      │  Gateway   │    │   Indexer   │    │  TxService  │
      │   (Go)     │    │    (Go)     │    │ (Java/Boot) │
      └─────┬──────┘    └──────┬──────┘    └──────┬──────┘
            │                  │                  │
         Redis             PostgreSQL          PostgreSQL
```

**Why three separate services?**
Each service handles a distinct failure domain. Gateway is stateless and horizontally scalable. Indexer is a single-writer to avoid data conflicts. TxService manages nonces which require strict serialization. Coupling them would mean a cache bug takes down transaction sending, or an indexer crash blocks API requests.

**Why Hoodi testnet, not a local devnet?**
Reorgs, variable block times, and unreliable RPCs only happen on a real network. A local Hardhat/Anvil node always returns instantly, never reorgs, and never drops connections — it would validate nothing about the problems this project exists to solve.

---

## Gateway

### Decision: Finality-aware cache TTLs

| Block reference | TTL | Why |
|---|---|---|
| `finalized` tag | 5 min | Finalized data cannot be reverted by reorg. Long cache is safe. |
| `safe` tag | 10 s | Safe but not finalized — extremely unlikely to revert, but the spec reserves the right. Short-to-medium hedge. |
| `latest` tag | 1 s | Changes every block (~12s). 1s gives a hit for concurrent requests within the same block without serving stale data across blocks. |
| `pending` tag | never | Pending state is speculative. Caching it would serve stale mempool data. |
| `earliest` tag | 5 min | Genesis block never changes. |
| Hex block number ≤ finalized | 5 min | Once finalized, block data is immutable. Same reasoning as `finalized` tag. |
| Hex block number > finalized, ≤ latest | 1 s | Exists on-chain but could be reorged. Same treatment as `latest`. |
| Hex block number > latest | never | Block doesn't exist yet. Caching null would poison future requests. |
| `{blockHash: "0x..."}` | 5 min | A specific block hash either exists or doesn't. If it exists, its data is immutable (hash commits to content). |
| Static methods (eth_chainId, net_version) | 5 min | These never change for a given network. |
| Gas methods (eth_gasPrice, etc.) | 1 s | Change every block, but concurrent requests within ~1s benefit from dedup. |

**Why not cache null results?**
Some methods return `null` to mean "not available yet" (e.g., a receipt for a pending tx). Caching null would mean subsequent callers keep seeing "not found" even after the data appears. We only cache non-null successful responses.

**Why not cache error responses?**
An upstream error is often transient (timeout, rate limit). Caching it would amplify a brief outage into a prolonged one.

### Decision: Canonicalize params for cache key only, not for upstream

The gateway normalizes params (lowercase hex, strip leading zeros, sort JSON keys) only when computing the cache key. The original raw params are forwarded to the upstream unchanged.

**Why?** The gateway is a proxy. Rewriting params risks breaking clients that depend on specific formatting, and makes debugging harder (the upstream sees different params than the client sent). Canonicalization errors only affect cache hit rate, not correctness.

**Address-aware normalization:** Ethereum addresses are always 40 hex characters. `canonicalizeHex` lowercases them but does not strip leading zeros, because `0x000...abc` and `0x abc` are different addresses.

### Decision: Multi-upstream with health scoring

Upstreams are scored on three dimensions: block height lag, error rate, and latency. The selector picks the highest-scoring upstream and falls back to the next on failure.

**Why not round-robin?** A lagging or erroring node would serve stale or failed responses to 1/N of all requests. Health scoring routes traffic away from degraded nodes automatically.

**Why still try unhealthy nodes?** If all nodes score 0, the selector still tries the least-bad one. A degraded response is better than no response. The health checker periodically re-probes, so a recovered node regains score over time.

### Decision: Rate limiting by API key with Redis

Rate limits use Redis so multiple gateway instances share the same counters.

**Why fail-open on Redis errors?** A Redis outage is temporary. Rejecting all traffic during a Redis blip is worse than briefly allowing unthrottled access. The rate limiter protects against sustained abuse, not against a burst during a 5-second Redis restart.

---

## Indexer

### Decision: Single-block-at-a-time sequential scan

The indexer processes blocks one at a time, in order, never skipping. After processing block N, it writes a checkpoint (block number + block hash) and moves to N+1.

**Why not batch or parallelize?** Reorg detection requires checking that block N+1's `parentHash` matches block N's hash. Parallel processing would break this chain. Batching would mean a failure mid-batch requires figuring out which blocks succeeded — single-block granularity makes recovery trivial: restart from the checkpoint.

### Decision: parentHash chain validation for reorg detection

Before processing any block, the indexer verifies:

```
chain[N+1].parentHash == checkpoint.blockHash
```

If this fails, a reorg has occurred. The indexer walks backwards to find the fork point (the latest block where chain and DB agree), rolls back everything after it, and resumes forward scanning.

**Why parentHash instead of polling `eth_getBlockByNumber` for changes?** ParentHash verification is local (one comparison) and catches reorgs immediately at the point of processing. Polling-based detection has a lag window where stale data has already been written.

**Why walk backwards instead of a binary search for the fork point?** Most reorgs are 1-2 blocks deep. Linear backward scan is O(1) in the common case. Binary search adds complexity for a scenario (deep reorg) that's extremely rare on PoS chains.

### Decision: Idempotent writes with `INSERT ON CONFLICT DO NOTHING`

Each transfer record's unique key is `(tx_hash, log_index)`. Re-processing the same block (e.g., after a crash before the checkpoint update) silently skips already-written records.

**Why not delete-and-rewrite?** Delete-rewrite creates a window where data is missing. `ON CONFLICT DO NOTHING` means the data is always either absent (not yet written) or correct (written once). No intermediate state.

### Decision: All writes in a single database transaction per block

Transfers, balance updates, and checkpoint — all in one transaction. Either the entire block is persisted or nothing is.

**Why?** Without this, a crash between writing transfers and updating the checkpoint would cause the indexer to re-process the block, potentially double-counting balances. The idempotent `ON CONFLICT` handles duplicate transfers, but balance updates (`amount = amount + value`) are not idempotent. The transaction ensures balance and transfer consistency.

### Decision: Track confirmation level per transfer

Each transfer record stores a `confirmation_level` (latest / safe / finalized). The indexer updates this as blocks progress through the confirmation pipeline.

**Why?** Downstream consumers (e.g., the exchange's deposit crediting system) need to know whether a deposit is final. Crediting a user's balance on a `latest` confirmation and having it reorged away means the exchange loses money. Exposing the level lets consumers make their own risk decisions.

---

## TxService

### Decision: Database row lock for nonce allocation

Nonce allocation uses `SELECT ... FOR UPDATE` on a per-address counter row. This serializes all nonce allocations for the same sender address.

**Why database lock instead of in-memory counter?** Multiple TxService instances may run concurrently. An in-memory counter would allow two instances to allocate the same nonce. The database lock is the coordination point.

**Why not derive nonce from `eth_getTransactionCount`?** The on-chain count only reflects confirmed transactions. If you've sent 5 pending transactions, `getTransactionCount` still returns the old value. You'd allocate duplicate nonces. The database counter tracks what has been *allocated*, not what has been *confirmed*.

**Why never recycle nonces?** A "failed" allocation (e.g., the broadcast times out) may actually have succeeded — the transaction could be in some node's mempool. Reusing that nonce would create a replacement transaction. We burn the nonce and move forward. Nonce gaps are resolved by the stuck transaction replacement logic.

### Decision: Same-nonce gas price replacement for stuck transactions

A transaction that's been BROADCAST for too long (>30 blocks, ~6 minutes) gets replaced: same nonce, same payload, but 10%+ higher gas price.

**Why 10% minimum?** EIP-1559 requires replacement transactions to have at least 10% higher `maxPriorityFeePerGas`. Anything less is rejected by nodes with "replacement transaction underpriced".

**Why not just wait longer?** On a congested network, a low-gas transaction may never be included. The mempool can evict it after hours, and the nonce is stuck — all subsequent transactions are blocked. Active replacement is necessary for liveness.

**Why cap the maximum gas price?** Without a cap, repeated replacements could escalate gas cost to an absurd level. Better to mark as FAILED and alert an operator than to spend $1000 in gas on a $10 transfer.

### Decision: State machine with explicit transitions

```
PENDING → BROADCAST → INCLUDED → FINALIZED
                ↘                 ↗
                 → REPLACED →  (new tx follows same path)
                ↘
                 → FAILED
```

**Why not just PENDING / SUCCESS / FAILED?** The intermediate states are operationally critical:

- **BROADCAST**: The exchange knows the tx was sent. If TxService crashes, it knows what to track when it restarts.
- **INCLUDED**: The tx is on-chain but could reorg. The exchange can show "confirming..." to the user.
- **FINALIZED**: Safe to update the user's off-chain balance. This is the only state where the exchange should consider the withdrawal complete.
- **REPLACED**: Audit trail. Shows that the original tx was stuck and replaced, preserving the history of what happened.

**Why allow INCLUDED → BROADCAST rollback?** If a reorg removes the block containing our transaction, the receipt disappears. The transaction may still be in the mempool and could be re-included in a different block. Rolling back to BROADCAST lets the tracker pick it up again naturally.

### Decision: Idempotent withdrawal creation

The withdrawal API uses an `idempotency_id` (client-provided). Duplicate requests return the existing withdrawal instead of creating a new one.

**Why?** Network issues cause retries. Without idempotency, a timeout + retry would create two withdrawals and send double the amount. The `idempotency_id` is a UNIQUE constraint — the database enforces exactly-once creation.

---

## Cross-cutting decisions

### Why Go for Gateway and Indexer, Java for TxService?

Pragmatic choice. Go's lightweight concurrency model (goroutines) fits the gateway's many-connections-few-CPU-work pattern and the indexer's single-loop-with-HTTP-calls pattern. Java/Spring Boot provides mature transaction management (JPA, `@Transactional`) that simplifies the TxService's database-heavy nonce and state machine logic. Also demonstrates proficiency in both ecosystems.

### Why Redis for cache and rate limiting, PostgreSQL for persistent state?

Cache and rate limit counters are ephemeral — losing them means a brief cache miss or a brief rate limit reset, not data loss. Redis is purpose-built for this: fast, TTL-native, atomic operations.

Transfer records, balances, nonces, and withdrawal state must survive restarts and be transactionally consistent. PostgreSQL provides ACID transactions, which are essential for the idempotent write and nonce allocation patterns.

### Why no message queue between services?

The three services interact only through the blockchain and their own databases. The gateway doesn't talk to the indexer; the indexer doesn't talk to TxService. This is by design — they are independently deployable and independently recoverable. A message queue would introduce a shared dependency and ordering concerns that don't exist in the current architecture.

### Why testnet only, with all secrets from environment variables?

This is a portfolio project. Running on mainnet would risk real funds for no benefit. But the code follows mainnet-grade practices (env-only secrets, no hardcoded keys, `.env` in `.gitignore`) so the patterns transfer directly.
