// Package client provides utilities to interact with external APIs such as Etherscan, Binance, etc
package client

import (
	"math/big"
	"time"
)

// KlineData is the return type of PriceClient GetETHUSDT
type KlineData struct {
	ClosePrice float64
}

// TransactionData represents the simplified transaction result from the API calls
type TransactionData struct {
	BlockNumber uint64
	Hash        string
	GasUsed     uint64
	GasPriceWei *big.Int
	Timestamp   time.Time
}

// PriceClient defines the interface for fetching price data. Mostly for dependency injection
type PriceClient interface {
	GetETHUSDT(timestamp time.Time) (*KlineData, error)
}

// TransactionClient defines the interface from fetching transactions data from the client
type TransactionClient interface {
	GetTransactionReceipt(hash string) (*TransactionData, error)
	GetLatestTransaction() (*TransactionData, error)
	ListTransactions(offset *int, startBlock *uint64, endBlock *uint64, page *int) ([]TransactionData, error)
}
