CREATE TABLE IF NOT EXISTS transactions (
    transaction_hash     TEXT PRIMARY KEY,
    block_number         BIGINT NOT NULL,
    timestamp            TIMESTAMPTZ NOT NULL,
    gas_used             BIGINT NOT NULL,
    gas_price_wei        NUMERIC NOT NULL,
    transaction_fee_eth  NUMERIC,
    transaction_fee_usdt NUMERIC,
    eth_usdt_price       NUMERIC
);

