package domain

import (
	"context"
	"time"

	"github.com/winQe/uniswap-fee-tracker/internal/types"
)

// PriceManagerInterface defines interface for price manager
type PriceManagerInterface interface {
	GetETHUSDT(timestamp time.Time) (float64, error)
}

// TransactionManagerInterface defines interface for transaction manager
type TransactionManagerInterface interface {
	GetLatestBlockNumber() (uint64, error)
	GetTransaction(hash string) (*types.TxWithPrice, error)
	BatchProcessTransactions(startBlock uint64, endBlock uint64, ctx context.Context) ([]types.TxWithPrice, error)
	BatchProcessTransactionsByTimestamp(startTime time.Time, endTime time.Time, ctx context.Context) ([]types.TxWithPrice, error)
}
