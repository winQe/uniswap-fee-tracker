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

// BatchProcessTransactions fetches and processes transactions within the given block range.
// It utilizes concurrent worker to fetch and process transactions
func (tm *TransactionManager) BatchProcessTransactions(startBlock uint64, endBlock uint64, ctx context.Context) ([]TxWithPrice, error) {
	var allTransactions []TxWithPrice

	batchSize := 100
	numWorkers := 10

	pages := make(chan int)
	results := make(chan TxWithPrice)
	stopSignal := make(chan struct{})

	var wg sync.WaitGroup

	// Ensure stopSignal is closed only once
	var once sync.Once

	// Worker function: fetches and processes transactions from pages channel
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-stopSignal:
				return
			case page, ok := <-pages:
				if !ok {
					return
				}
				// Fetch transactions for the current page
				transactions, err := tm.transactionClient.ListTransactions(&batchSize, &startBlock, nil, &page)
				if err != nil {
					fmt.Printf("Error fetching page %d: %v\n", page, err)
					continue
				}

				// Process transactions
				for _, tx := range transactions {
					txWithPrice, err := tm.processTransaction(tx)
					if err != nil {
						fmt.Printf("Error processing transaction %s: %v\n", tx.Hash, err)
						continue
					}
					select {
					case results <- *txWithPrice:
					case <-ctx.Done():
						return
					}
				}

				// If fewer transactions than batchSize are returned, it's likely the last page
				if len(transactions) < batchSize {
					// Signal to stop further page dispatching
					once.Do(func() {
						close(stopSignal)
					})
				}
			}
		}
	}

	// Start worker goroutines
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Dispatcher: sends page numbers to pages channel
	go func() {
		defer close(pages)
		page := 1
		for {
			select {
			case <-ctx.Done():
				return
			case <-stopSignal:
				return
			default:
				select {
				case pages <- page:
					page++
				case <-ctx.Done():
					return
				case <-stopSignal:
					return
				}
			}
		}
	}()

	// Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect processed results
	for tx := range results {
		allTransactions = append(allTransactions, tx)
	}

	return allTransactions, nil
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
