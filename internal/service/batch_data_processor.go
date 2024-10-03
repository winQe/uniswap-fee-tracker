package service

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

// BatchDataProcessor defines the interface for processing batch jobs with GoRoutines
type BatchDataProcessor interface {
	ProcessBatchJob(jobID string, startTime, endTime int64) error
}

// BatchDataProcessorImpl is the concrete implementation of BatchDataProcessor.
type BatchDataProcessorImpl struct {
	txDbQuery db.Querier
	jobCache  cache.JobsStore
	txManager domain.TransactionManagerInterface
}

// NewBatchDataProcessor initializes a new BatchDataProcessorImpl.
func NewBatchDataProcessor(txDbQuery db.Querier, jobCache cache.JobsStore, txManager domain.TransactionManagerInterface) *BatchDataProcessorImpl {
	return &BatchDataProcessorImpl{
		txDbQuery: txDbQuery,
		jobCache:  jobCache,
		txManager: txManager,
	}
}

// ProcessBatchJob processes the batch job asynchronously.
func (bdp *BatchDataProcessorImpl) ProcessBatchJob(jobID string, startTime, endTime int64) error {
	// Update job status to 'running'
	bdp.updateJobStatus(jobID, "running", "")

	// Create a new context for the batch processing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Convert unix time to time.Time
	startTs := time.Unix(startTime, 0)
	endTs := time.Unix(endTime, 0)

	// Execute the batch processing
	result, err := bdp.txManager.BatchProcessTransactionsByTimestamp(startTs, endTs, ctx)
	for _, tx := range result {
		err := bdp.txDbQuery.InsertTransaction(context.Background(), db.InsertTransactionParams{
			TransactionHash:    tx.Hash,
			BlockNumber:        int64(tx.BlockNumber),
			Timestamp:          tx.Timestamp,
			GasUsed:            int64(tx.GasUsed),
			GasPriceWei:        tx.GasPriceWei.Int64(),
			TransactionFeeEth:  pgtype.Float8{Float64: tx.TransactionFeeETH, Valid: true},
			TransactionFeeUsdt: pgtype.Float8{Float64: tx.TransactionFeeUSDT, Valid: true},
			EthUsdtPrice:       pgtype.Float8{Float64: tx.ETHUSDTPrice, Valid: true},
		})
		if err != nil {
			log.Printf("Error inserting transaction %s into DB: %v\n", tx.Hash, err)
			continue
		}
	}

	if err != nil {
		// Update job status to 'failed' with error message
		bdp.updateJobStatus(jobID, "failed", err.Error())
		return err
	}
	// Update job status to 'completed' with a success message
	if err := bdp.updateJobStatus(jobID, "completed", "Batch job completed successfully."); err != nil {
		return err
	}

	return nil
}

// updateJobStatus updates the status and result of a batch job in Redis.
func (bdp *BatchDataProcessorImpl) updateJobStatus(jobID, status, result string) error {
	// Retrieve the current job data
	jobData, err := bdp.jobCache.GetJob(jobID)
	if err != nil {
		log.Printf("Failed to retrieve job %s for status update: %v", jobID, err)
		return err
	}

	var job cache.BatchJobResponse
	if err := utils.DeserializeFromJSON(jobData, &job); err != nil {
		log.Printf("Failed to deserialize job %s data: %v", jobID, err)
		return err
	}

	// Update status and result
	job.Status = status
	job.UpdatedAt = time.Now().Unix()
	if result != "" {
		job.Result = result
	}

	// Serialize updated job
	updatedJobData, err := utils.SerializeToJSON(job)
	if err != nil {
		log.Printf("Failed to serialize updated job %s data: %v", jobID, err)
		return err
	}

	// Store the updated job data back in Redis
	err = bdp.jobCache.SetJob(job.ID, updatedJobData)
	if err != nil {
		log.Printf("Failed to update job %s in Redis: %v", job.ID, err)
		return err
	}

	return nil
}
