package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
)

// LiveDataRecorder records
type LiveDataRecorder struct {
	lastBlockNumber    uint64
	transactionManager domain.TransactionManagerInterface
	dbQuerier          db.Querier
}

// NewLiveDataRecorder initializes a new LiveDataRecorder instance.
func NewLiveDataRecorder(dbQuerier db.Querier, transactionManager domain.TransactionManagerInterface) *LiveDataRecorder {
	lastBlockNumber, err := transactionManager.GetLatestBlockNumber()
	if err != nil {
		log.Fatalf("Failed to get the latest block number: %v\n", err)
	}
	return &LiveDataRecorder{
		lastBlockNumber:    lastBlockNumber,
		transactionManager: transactionManager,
		dbQuerier:          dbQuerier,
	}
}

// Run starts the Recorder to execute tasks every 60 seconds.
// It listens for context cancellation to gracefully shut down.
func (ldr *LiveDataRecorder) Run(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Println("LiveDataRecorder started.")

	for {
		select {
		case <-ctx.Done():
			log.Println("LiveDataRecorder shutting down.")
			return
		case <-ticker.C:
			ldr.recordNewTransactions()
		}
	}
}

// executeTask contains the logic to fetch and process new transactions.
func (ldr *LiveDataRecorder) recordNewTransactions() {
	log.Println("Fetching and processing transactions.")

	latestBlock, err := ldr.transactionManager.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Error fetching latest block number: %v\n", err)
		return
	}

	// Fetch and process transactions from lastBlockNumber+1 to latestBlock.
	if latestBlock > ldr.lastBlockNumber {
		startBlock := ldr.lastBlockNumber + 1
		endBlock := latestBlock

		transactions, err := ldr.transactionManager.BatchProcessTransactions(startBlock, endBlock, context.Background())
		if err != nil {
			log.Printf("Error processing transactions from block %d to %d: %v\n", startBlock, endBlock, err)
			return
		}

		for _, tx := range transactions {
			// Example: Insert transaction into the database.
			err := ldr.dbQuerier.InsertTransaction(context.Background(), db.InsertTransactionParams{
				TransactionHash:    tx.Hash,
				BlockNumber:        int64(tx.BlockNumber),
				Timestamp:          tx.Timestamp,
				GasUsed:            int64(tx.GasUsed),
				GasPriceWei:        tx.GasPriceWei.Int64(),
				TransactionFeeEth:  sql.NullFloat64{Float64: tx.TransactionFeeETH, Valid: true},
				TransactionFeeUsdt: sql.NullFloat64{Float64: tx.TransactionFeeUSDT, Valid: true},
				EthUsdtPrice:       sql.NullFloat64{Float64: tx.ETHUSDTPrice, Valid: true},
			})
			if err != nil {
				log.Printf("Error inserting transaction %s into DB: %v\n", tx.Hash, err)
				continue
			}
		}

		// Update the last processed block number.
		ldr.lastBlockNumber = endBlock
		numTxProcessed := len(transactions)
		log.Printf("Processed %d transactions up to block %d.\n", numTxProcessed, endBlock)
	} else {
		log.Println("No new transactions to process.")
	}
}
