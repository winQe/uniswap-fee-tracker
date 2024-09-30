CREATE TABLE transactions (
    transaction_hash     TEXT PRIMARY KEY,
    block_number         BIGINT NOT NULL,
    timestamp            TIMESTAMPTZ NOT NULL,
    gas_used             BIGINT NOT NULL,
    gas_price_wei        NUMERIC NOT NULL,
    transaction_fee_eth  NUMERIC, -- Calculated as gas_used * gas_price_wei / 1e18
    transaction_fee_usdt NUMERIC, -- Calculated as transaction_fee_eth * eth_usdt_price
    eth_usdt_price       NUMERIC  -- ETH/USDT price at transaction time
);

