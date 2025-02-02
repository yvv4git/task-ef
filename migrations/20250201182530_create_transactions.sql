-- +goose Up
-- +goose StatementBegin
CREATE TABLE eth_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_hash TEXT,
    block_number BIGINT,
    from_address TEXT,
    gas BIGINT,
    gas_price NUMERIC,
    hash TEXT UNIQUE,
    input TEXT,
    nonce BIGINT,
    to_address TEXT,
    transaction_index INTEGER,
    value NUMERIC,
    type INTEGER,
    v TEXT,
    r TEXT,
    s TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE eth_transactions;
-- +goose StatementEnd