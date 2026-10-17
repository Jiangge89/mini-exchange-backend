CREATE TABLE IF NOT EXISTS transfers (
    id              BIGSERIAL PRIMARY KEY,
    block_number    BIGINT      NOT NULL,
    block_hash      BYTEA       NOT NULL,
    tx_hash         BYTEA       NOT NULL,
    log_index       INTEGER     NOT NULL,
    from_addr       BYTEA       NOT NULL,
    to_addr         BYTEA       NOT NULL,
    value           NUMERIC     NOT NULL,
    confirmation_level TEXT     NOT NULL DEFAULT 'latest',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (tx_hash, log_index)
);

CREATE INDEX idx_transfers_from ON transfers (from_addr, block_number DESC);
CREATE INDEX idx_transfers_to   ON transfers (to_addr, block_number DESC);
CREATE INDEX idx_transfers_block ON transfers (block_number);

CREATE TABLE IF NOT EXISTS balances (
    address BYTEA   PRIMARY KEY,
    amount  NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS checkpoints (
    id           INTEGER PRIMARY KEY DEFAULT 1,
    block_number BIGINT  NOT NULL,
    block_hash   BYTEA   NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT single_row CHECK (id = 1)
);
