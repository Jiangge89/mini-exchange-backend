CREATE TABLE withdrawals (
    id              BIGSERIAL   PRIMARY KEY,
    idempotency_id  VARCHAR(255) NOT NULL UNIQUE,
    to_address      VARCHAR(42) NOT NULL,
    amount          NUMERIC(78, 0) NOT NULL,
    nonce           INTEGER,
    tx_hash         VARCHAR(66),
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    gas_price       NUMERIC(78, 0),
    block_number    BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withdrawals_status ON withdrawals (status);
CREATE INDEX idx_withdrawals_idempotency ON withdrawals (idempotency_id);

CREATE TABLE nonce_counter (
    address     VARCHAR(42) PRIMARY KEY,
    next_nonce  INTEGER NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
