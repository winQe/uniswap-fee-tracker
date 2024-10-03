-- name: InsertTransaction :exec
INSERT INTO transactions (
    transaction_hash,
    block_number,
    timestamp,
    gas_used,
    gas_price_wei,
    transaction_fee_eth,
    transaction_fee_usdt,
    eth_usdt_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetTransactionByHash :one
SELECT
    transaction_hash,
    block_number,
    timestamp,
    gas_used,
    gas_price_wei,
    transaction_fee_eth,
    transaction_fee_usdt,
    eth_usdt_price
FROM transactions
WHERE transaction_hash = $1;

-- name: GetTransactionsByBlockNumber :many
SELECT
    transaction_hash,
    block_number,
    timestamp,
    gas_used,
    gas_price_wei,
    transaction_fee_eth,
    transaction_fee_usdt,
    eth_usdt_price
FROM transactions
WHERE block_number = $1
ORDER BY timestamp DESC;

-- name: GetTransactionsByTimeRange :many
SELECT
    transaction_hash,
    block_number,
    timestamp,
    gas_used,
    gas_price_wei,
    transaction_fee_eth,
    transaction_fee_usdt,
    eth_usdt_price
FROM transactions
WHERE timestamp BETWEEN $1 AND $2
ORDER BY timestamp DESC;

-- name: GetLatestTransactions :many
SELECT *
FROM transactions
ORDER BY timestamp DESC
LIMIT $1;
