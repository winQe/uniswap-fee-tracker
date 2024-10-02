// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Transactions struct {
	TransactionHash    string        `json:"transaction_hash"`
	BlockNumber        int64         `json:"block_number"`
	Timestamp          time.Time     `json:"timestamp"`
	GasUsed            int64         `json:"gas_used"`
	GasPriceWei        int64         `json:"gas_price_wei"`
	TransactionFeeEth  pgtype.Float8 `json:"transaction_fee_eth"`
	TransactionFeeUsdt pgtype.Float8 `json:"transaction_fee_usdt"`
	EthUsdtPrice       pgtype.Float8 `json:"eth_usdt_price"`
}
