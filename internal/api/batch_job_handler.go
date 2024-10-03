package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
	"github.com/winQe/uniswap-fee-tracker/internal/service"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

// BatchJobHandler handles batch job-related HTTP requests.
type BatchJobHandler struct {
	txDbQuery          db.Querier
	jobCache           cache.JobsStore
	txManager          domain.TransactionManagerInterface
	batchDataProcessor service.BatchDataProcessor
}

// NewBatchJobHandler initializes a new BatchJobHandler with the given dependencies.
func NewBatchJobHandler(txDbQuery db.Querier, jobCache cache.JobsStore, txManager domain.TransactionManagerInterface, batchDataProcessor service.BatchDataProcessor) *BatchJobHandler {
	return &BatchJobHandler{
		txDbQuery:          txDbQuery,
		jobCache:           jobCache,
		txManager:          txManager,
		batchDataProcessor: batchDataProcessor,
	}
}

// CreateBatchJob schedules a new batch job for historical data recording.
// It accepts 'start_time' and 'end_time' as query parameters in Unix epoch seconds.
func (bh *BatchJobHandler) CreateBatchJob(ctx *gin.Context) {
	// Parse query parameters
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Missing 'start_time' or 'end_time' query parameters"})
		return
	}

	// Convert query parameters to int64
	startTime, err := utils.ParseUnixTime(startTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid 'start_time' format. Must be Unix epoch time in seconds."})
		return
	}

	endTime, err := utils.ParseUnixTime(endTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid 'end_time' format. Must be Unix epoch time in seconds."})
		return
	}

	if endTime <= startTime {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "End time must be after start time"})
		return
	}

	// Generate a unique ID for the batch job
	jobID := uuid.New().String()
	currentTime := time.Now().Unix()

	// Create BatchJobResponse with initial status 'pending'
	job := cache.BatchJob{
		ID:        jobID,
		Status:    "pending",
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Result:    "",
	}

	// Serialize job to JSON
	jobData, err := utils.SerializeToJSON(job)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to serialize batch job data"})
		return
	}

	// Store the batch job in Redis with status 'pending'
	err = bh.jobCache.SetJob(jobID, jobData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to store batch job in Redis"})
		return
	}

	// Run batch job in the background with Goroutines
	go bh.batchDataProcessor.ProcessBatchJob(jobID, startTime, endTime)

	ctx.JSON(http.StatusCreated, job)
}

// GetBatchJob retrieves the status and details of a specific batch job.
func (bh *BatchJobHandler) GetBatchJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if !utils.IsValidUUID(jobID) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid batch job ID format"})
		return
	}

	// Retrieve job data from Redis
	jobData, err := bh.jobCache.GetJob(jobID)
	if err != nil {
		if err == cache.ErrJobNotFound {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse{Error: "Batch job not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to retrieve batch job"})
		return
	}

	var job cache.BatchJob
	if err := utils.DeserializeFromJSON(jobData, &job); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to parse batch job data"})
		return
	}

	ctx.JSON(http.StatusOK, job)
}

// ListBatchJobs lists all batch jobs without pagination.
func (bh *BatchJobHandler) ListBatchJobs(ctx *gin.Context) {
	// Fetch all batch jobs from Redis
	allJobs, err := bh.jobCache.GetAllJobs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to retrieve batch jobs"})
		return
	}

	var filteredJobs []cache.BatchJob
	for _, jobData := range allJobs {
		var job cache.BatchJob
		if err := utils.DeserializeFromJSON(jobData, &job); err != nil {
			continue // Skip malformed job entries
		}

		// Optionally filter by status if provided
		status := ctx.Query("status")
		if status != "" && job.Status != status {
			continue
		}

		filteredJobs = append(filteredJobs, job)
	}

	ctx.JSON(http.StatusOK, filteredJobs)
}
