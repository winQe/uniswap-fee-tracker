package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	"github.com/winQe/uniswap-fee-tracker/internal/mocks"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

func TestCreateBatchJob_Success(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create mock instances
	mockJobsStore := new(mocks.MockJobsStore)
	mockBatchDataProcessor := mocks.NewMockBatchDataProcessor() // Use constructor
	mockTxDbQuery := new(mocks.MockQuerier)
	mockTxManager := new(mocks.MockTransactionManager)

	// Initialize BatchJobHandler with mocks
	handler := NewBatchJobHandler(mockTxDbQuery, mockJobsStore, mockTxManager, mockBatchDataProcessor)

	// Create a test router and register the handler
	router := gin.Default()
	router.POST("/batch-jobs", handler.CreateBatchJob)

	// Define test input
	startTime := time.Now().Unix() - 3600 // 1 hour ago
	endTime := time.Now().Unix()

	// Create expected fields (excluding ID and timestamps)
	expectedStatus := "pending"
	expectedResult := ""

	// Set up mock expectations
	mockJobsStore.On("SetJob", mock.AnythingOfType("string"), mock.MatchedBy(func(data []byte) bool {
		var job cache.BatchJob
		err := json.Unmarshal(data, &job)
		if err != nil {
			return false
		}
		// Validate fields (excluding ID and timestamps)
		return job.Status == expectedStatus &&
			job.StartTime == startTime &&
			job.EndTime == endTime &&
			job.Result == expectedResult
	})).Return(nil)

	mockBatchDataProcessor.On("ProcessBatchJob", mock.AnythingOfType("string"), startTime, endTime).Return(nil)

	// Create a test HTTP request with query parameters
	req, err := http.NewRequest("POST", "/batch-jobs", nil)
	assert.NoError(t, err)

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_time", strconv.FormatInt(startTime, 10))
	q.Add("end_time", strconv.FormatInt(endTime, 10))
	req.URL.RawQuery = q.Encode()

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Wait for ProcessBatchJob to be called or timeout after 1 second
	select {
	case <-mockBatchDataProcessor.CalledChan:
		// Method was called
	case <-time.After(1 * time.Second):
		t.Fatal("ProcessBatchJob was not called within timeout")
	}

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Parse the response body
	var responseBody cache.BatchJob
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Assert the response body fields (excluding ID and timestamps)
	assert.Equal(t, expectedStatus, responseBody.Status)
	assert.Equal(t, startTime, responseBody.StartTime)
	assert.Equal(t, endTime, responseBody.EndTime)
	assert.Equal(t, expectedResult, responseBody.Result)
	assert.NotEmpty(t, responseBody.ID, "Job ID should not be empty")

	// Optionally, validate that the ID is a valid UUID
	_, err = uuid.Parse(responseBody.ID)
	assert.NoError(t, err, "Job ID should be a valid UUID")

	// Assert that SetJob was called with the correct parameters
	mockJobsStore.AssertCalled(t, "SetJob", mock.AnythingOfType("string"), mock.MatchedBy(func(data []byte) bool {
		var job cache.BatchJob
		err := json.Unmarshal(data, &job)
		if err != nil {
			return false
		}
		return job.Status == expectedStatus &&
			job.StartTime == startTime &&
			job.EndTime == endTime &&
			job.Result == expectedResult
	}))

	// Assert that ProcessBatchJob was called
	mockBatchDataProcessor.AssertCalled(t, "ProcessBatchJob", mock.AnythingOfType("string"), startTime, endTime)
}

func TestCreateBatchJob_SetJobFailure(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create mock instances
	mockJobsStore := new(mocks.MockJobsStore)
	mockBatchDataProcessor := new(mocks.MockBatchDataProcessor)
	mockTxDbQuery := new(mocks.MockQuerier)
	mockTxManager := new(mocks.MockTransactionManager)

	// Initialize BatchJobHandler with mocks
	handler := NewBatchJobHandler(mockTxDbQuery, mockJobsStore, mockTxManager, mockBatchDataProcessor)

	// Create a test router and register the handler
	router := gin.Default()
	router.POST("/batch-jobs", handler.CreateBatchJob)

	// Define test input
	startTime := time.Now().Unix() - 3600 // 1 hour ago
	endTime := time.Now().Unix()

	// Set up mock expectations for SetJob to return an error
	mockJobsStore.On("SetJob", mock.AnythingOfType("string"), mock.MatchedBy(func(data []byte) bool {
		var job cache.BatchJob
		err := json.Unmarshal(data, &job)
		if err != nil {
			return false
		}
		// Validate fields (excluding ID and timestamps)
		return job.Status == "pending" &&
			job.StartTime == startTime &&
			job.EndTime == endTime &&
			job.Result == ""
	})).Return(errors.New("Redis error"))

	// Create a test HTTP request with query parameters
	req, err := http.NewRequest("POST", "/batch-jobs", nil)
	assert.NoError(t, err)

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_time", strconv.FormatInt(startTime, 10))
	q.Add("end_time", strconv.FormatInt(endTime, 10))
	req.URL.RawQuery = q.Encode()

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// Parse the response body
	var responseBody ErrorResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Assert the error message
	assert.Equal(t, "Failed to store batch job in Redis", responseBody.Error)

	// Assert that SetJob was called with the correct parameters
	mockJobsStore.AssertCalled(t, "SetJob", mock.AnythingOfType("string"), mock.MatchedBy(func(data []byte) bool {
		var job cache.BatchJob
		err := json.Unmarshal(data, &job)
		if err != nil {
			return false
		}
		return job.Status == "pending" &&
			job.StartTime == startTime &&
			job.EndTime == endTime &&
			job.Result == ""
	}))

	// Assert that ProcessBatchJob was not called
	mockBatchDataProcessor.AssertNotCalled(t, "ProcessBatchJob", mock.Anything, mock.Anything, mock.Anything)
}
