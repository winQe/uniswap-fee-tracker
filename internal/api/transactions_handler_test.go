package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/mocks"
)

// TestGetTransactionHash tests the retrieval of a transaction by its hash.
func TestGetTransactionHash(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock Querier
	mockQuerier := new(mocks.MockQuerier)

	// Sample transaction data
	sampleTx := db.Transactions{
		TransactionHash:    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		BlockNumber:        123456,
		Timestamp:          time.Unix(1617181723, 0).UTC(),
		GasUsed:            21000,
		GasPriceWei:        1000000000,
		TransactionFeeEth:  pgtype.Float8{Float64: 0.021, Valid: true},
		TransactionFeeUsdt: pgtype.Float8{Float64: 42.0, Valid: true},
		EthUsdtPrice:       pgtype.Float8{Float64: 2000.0, Valid: true},
	}

	// Set up expectations
	mockQuerier.On("GetTransactionByHash", mock.Anything, sampleTx.TransactionHash).Return(sampleTx, nil)

	// Initialize TransactionHandler
	handler := NewTransactionHandler(mockQuerier)

	// Set up Gin router
	router := gin.Default()
	router.GET("/transactions/:hash", handler.getTransactionHash)

	// Create a test request
	req, _ := http.NewRequest("GET", "/transactions/"+sampleTx.TransactionHash, nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `{
		"transaction_hash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"block_number": 123456,
		"timestamp": 1617181723,
		"gas_used": 21000,
		"gas_price_wei": 1000000000,
		"transaction_fee_eth": 0.021,
		"transaction_fee_usdt": 42,
		"eth_usdt_price": 2000
	}`
	assert.JSONEq(t, expectedBody, resp.Body.String())

	// Assert that the expectations were met
	mockQuerier.AssertExpectations(t)
}

// TestGetLatestTransactions_DefaultLimit tests the retrieval of the latest transactions with the default limit.
func TestGetLatestTransactions_DefaultLimit(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock Querier
	mockQuerier := new(mocks.MockQuerier)

	// Sample transactions data
	sampleTxs := []db.Transactions{
		{
			TransactionHash:    "0xhash1",
			BlockNumber:        123456,
			Timestamp:          time.Unix(1617181723, 0).UTC(),
			GasUsed:            21000,
			GasPriceWei:        1000000000,
			TransactionFeeEth:  pgtype.Float8{Float64: 0.021, Valid: true},
			TransactionFeeUsdt: pgtype.Float8{Float64: 42.0, Valid: true},
			EthUsdtPrice:       pgtype.Float8{Float64: 2000.0, Valid: true},
		},
		{
			TransactionHash:    "0xhash2",
			BlockNumber:        123457,
			Timestamp:          time.Unix(1617181730, 0).UTC(),
			GasUsed:            22000,
			GasPriceWei:        1100000000,
			TransactionFeeEth:  pgtype.Float8{Float64: 0.022, Valid: true},
			TransactionFeeUsdt: pgtype.Float8{Float64: 44.0, Valid: true},
			EthUsdtPrice:       pgtype.Float8{Float64: 2000.0, Valid: true},
		},
		// Add more transactions as needed
	}

	// Set up expectations
	mockQuerier.On("GetLatestTransactions", mock.Anything, int32(10)).Return(sampleTxs, nil)

	// Initialize TransactionHandler
	handler := NewTransactionHandler(mockQuerier)

	// Set up Gin router
	router := gin.Default()
	router.GET("/transactions/latest", handler.getLatestTransactions)

	// Create a test request without 'limit'
	req, _ := http.NewRequest("GET", "/transactions/latest", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `[
		{
			"transaction_hash": "0xhash1",
			"block_number": 123456,
			"timestamp": 1617181723,
			"gas_used": 21000,
			"gas_price_wei": 1000000000,
			"transaction_fee_eth": 0.021,
			"transaction_fee_usdt": 42,
			"eth_usdt_price": 2000
		},
		{
			"transaction_hash": "0xhash2",
			"block_number": 123457,
			"timestamp": 1617181730,
			"gas_used": 22000,
			"gas_price_wei": 1100000000,
			"transaction_fee_eth": 0.022,
			"transaction_fee_usdt": 44,
			"eth_usdt_price": 2000
		}
	]`
	assert.JSONEq(t, expectedBody, resp.Body.String())

	// Assert that the expectations were met
	mockQuerier.AssertExpectations(t)
}

// TestGetTransactionByTimestamp tests the getTransactionByTimestamp handler.
func TestGetTransactionByTimestamp(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock Querier from the mocks package
	mockQuerier := new(mocks.MockQuerier)

	// Sample transactions data within the time range
	sampleTxs := []db.Transactions{
		{
			TransactionHash:    "0xhash3",
			BlockNumber:        123458,
			Timestamp:          time.Unix(1617181740, 0).UTC(), // 2021-03-31T12:09:00Z
			GasUsed:            21000,
			GasPriceWei:        1000000000,
			TransactionFeeEth:  pgtype.Float8{Float64: 0.021, Valid: true},
			TransactionFeeUsdt: pgtype.Float8{Float64: 420.0, Valid: true},
			EthUsdtPrice:       pgtype.Float8{Float64: 20000.0, Valid: true},
		},
		{
			TransactionHash:    "0xhash4",
			BlockNumber:        123459,
			Timestamp:          time.Unix(1617181750, 0).UTC(), // 2021-03-31T12:09:10Z
			GasUsed:            22000,
			GasPriceWei:        1100000000,
			TransactionFeeEth:  pgtype.Float8{Float64: 0.0242, Valid: true},
			TransactionFeeUsdt: pgtype.Float8{Float64: 484.0, Valid: true},
			EthUsdtPrice:       pgtype.Float8{Float64: 20000.0, Valid: true},
		},
	}

	// Define start and end times in Unix seconds
	startUnix := int64(1617181720) // 2021-03-31T12:08:40Z
	endUnix := int64(1617181760)   // 2021-03-31T12:09:20Z

	// Set up expectations for GetTransactionsByTimeRange
	params := db.GetTransactionsByTimeRangeParams{
		Timestamp:   time.Unix(startUnix, 0),
		Timestamp_2: time.Unix(endUnix, 0),
	}
	mockQuerier.On("GetTransactionsByTimeRange", mock.Anything, params).Return(sampleTxs, nil)

	// Initialize TransactionHandler with the mock Querier
	handler := NewTransactionHandler(mockQuerier)

	// Set up Gin router and register the route
	router := gin.Default()
	router.GET("/transactions", handler.getTransactionByTimestamp)

	// Create a test request with 'start' and 'end' query parameters in RFC3339 format
	startTimeStr := strconv.FormatInt(startUnix, 10)
	endTimeStr := strconv.FormatInt(endUnix, 10)
	reqURL := "/transactions?start=" + startTimeStr + "&end=" + endTimeStr
	req, _ := http.NewRequest("GET", reqURL, nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Define the expected JSON response
	expectedBody := `[
		{
			"transaction_hash": "0xhash3",
			"block_number": 123458,
			"timestamp": 1617181740,
			"gas_used": 21000,
			"gas_price_wei": 1000000000,
			"transaction_fee_eth": 0.021,
			"transaction_fee_usdt": 420.0,
			"eth_usdt_price": 20000.0
		},
		{
			"transaction_hash": "0xhash4",
			"block_number": 123459,
			"timestamp": 1617181750,
			"gas_used": 22000,
			"gas_price_wei": 1100000000,
			"transaction_fee_eth": 0.0242,
			"transaction_fee_usdt": 484.0,
			"eth_usdt_price": 20000.0
		}
	]`

	// Assert that the response body matches the expected JSON
	assert.JSONEq(t, expectedBody, resp.Body.String())

	// Assert that all expectations were met
	mockQuerier.AssertExpectations(t)
}
