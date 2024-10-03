// Package client provides utilities to interact with external APIs such as Etherscan, Binance, etc
package client

import (
	"time"

	"github.com/winQe/uniswap-fee-tracker/internal/types"
)

// KlineData is the return type of PriceClient GetETHUSDT
type KlineData struct {
	ClosePrice float64
}

// PriceClient defines the interface for fetching price data. Mostly for dependency injection
type PriceClient interface {
	GetETHUSDT(timestamp time.Time) (*KlineData, error)
}

// TransactionClient defines the interface from fetching transactions data from the client
type TransactionClient interface {
	GetTransactionReceipt(hash string) (*types.TransactionData, error)
	GetLatestTransaction() (*types.TransactionData, error)
	ListTransactions(offset *int, startBlock *uint64, endBlock *uint64, page *int) ([]types.TransactionData, error)
	GetBlockNumberByTimestamp(timestamp time.Time, before bool) (uint64, error)
}
