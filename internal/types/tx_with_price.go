package types

import (
	"math/big"
	"time"
)

// TransactionData represents the simplified transaction result from the API calls
type TransactionData struct {
	BlockNumber uint64
	Hash        string
	GasUsed     uint64
	GasPriceWei *big.Int
	Timestamp   time.Time
}

// TxWithPrice holds the processed transaction data
type TxWithPrice struct {
	TransactionData
	ETHUSDTPrice       float64
	TransactionFeeETH  float64
	TransactionFeeUSDT float64
}
