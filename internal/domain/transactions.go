package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/winQe/uniswap-fee-tracker/internal/client"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

// TransactionManager handles the business logic of computing the transaction price
type TransactionManager struct {
	transactionClient client.TransactionClient
	priceManager      PriceManagerInterface
}

type TxWithPrice struct {
	client.TransactionData
	ETHUSDTPrice       float64
	TransctionFeeETH   float64
	TransactionFeeUSDT float64
}

// GetLatestBlockNumber returns the latest transaction block number from the Uniswap V3 WETH-USDC pool
func (tm *TransactionManager) GetLatestBlockNumber() (uint64, error) {
	txData, err := tm.transactionClient.GetLatestTransaction()
	if err != nil {
		return 0, fmt.Errorf("failed to get the latest transaction from the API client: %v", err)
	}
	return txData.BlockNumber, nil
}

// GetTransaction queries transaction by hash and calculates its transaction price in USDT
// It utilizes concurrent workers to efficiently retrieve and handle transactions in batches.
func (tm *TransactionManager) GetTransaction(hash string) (*TxWithPrice, error) {
	txData, err := tm.transactionClient.GetTransactionReceipt(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash from API client: %v", err)
	}

	return tm.processTransaction(*txData)
}

// processTransaction fetches transaction receipt and calculates fees
func (tm *TransactionManager) processTransaction(tx client.TransactionData) (*TxWithPrice, error) {
	// Fetch transaction receipt
	txWithReceipt, err := tm.transactionClient.GetTransactionReceipt(tx.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Fetch ETH-USDT conversion rate at the transaction's timestamp
	ethUSDTConversionRate, err := tm.priceManager.GetETHUSDT(txWithReceipt.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to get ETH-USDT conversion rate: %v", err)
	}

	// Calculate fees
	feeETH := utils.ConvertToETH(txWithReceipt.GasPriceWei)
	feeUSDT := feeETH * ethUSDTConversionRate

	return &TxWithPrice{
		*txWithReceipt,
		ethUSDTConversionRate,
		feeETH,
		feeUSDT,
	}, nil
}
