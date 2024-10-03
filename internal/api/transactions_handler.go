package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

// TransactionResponse represents the JSON structure of a transaction in the API response.
type TransactionResponse struct {
	TransactionHash    string  `json:"transaction_hash"`
	BlockNumber        int64   `json:"block_number"`
	Timestamp          int64   `json:"timestamp"` // Unix epoch time in seconds
	GasUsed            int64   `json:"gas_used"`
	GasPriceWei        int64   `json:"gas_price_wei"`
	TransactionFeeEth  float64 `json:"transaction_fee_eth"`
	TransactionFeeUsdt float64 `json:"transaction_fee_usdt"`
	EthUsdtPrice       float64 `json:"eth_usdt_price"`
}

type TransactionHandler struct {
	txDbQuery db.Querier
}

// NewTransactionHandler initializes a new TransactionHandler with the given dependencies.
func NewTransactionHandler(txDbQuery db.Querier) *TransactionHandler {
	return &TransactionHandler{
		txDbQuery: txDbQuery,
	}
}

func (th *TransactionHandler) getTransactionHash(ctx *gin.Context) {
	txHash := ctx.Param("hash")
	txHash = utils.SanitizeTransactionHash(txHash)

	if txHash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing transaction hash"})
		return
	}

	// Fetch transaction from the database
	transaction, err := th.txDbQuery.GetTransactionByHash(ctx, txHash)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := TransactionResponse{
		TransactionHash:    transaction.TransactionHash,
		BlockNumber:        transaction.BlockNumber,
		Timestamp:          transaction.Timestamp.Unix(),
		GasUsed:            transaction.GasUsed,
		GasPriceWei:        transaction.GasPriceWei,
		TransactionFeeEth:  float64(transaction.TransactionFeeEth.Float64),
		TransactionFeeUsdt: float64(transaction.TransactionFeeUsdt.Float64),
		EthUsdtPrice:       float64(transaction.EthUsdtPrice.Float64),
	}

	ctx.JSON(http.StatusOK, response)
}

// getLatestTransactions retrieves the latest transactions.
// Query parameter 'limit' can be used to specify the number of transactions to retrieve.
// Defaults to 10 if not provided.
func (th *TransactionHandler) getLatestTransactions(ctx *gin.Context) {
	// Default limit
	limit := int32(10)

	// Override limit if provided
	if l, exists := ctx.GetQuery("limit"); exists {
		parsedLimit, err := strconv.ParseInt(l, 10, 32)
		if err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	transactions, err := th.txDbQuery.GetLatestTransactions(ctx, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var response []TransactionResponse
	for _, tx := range transactions {
		txResp := TransactionResponse{
			TransactionHash:    tx.TransactionHash,
			BlockNumber:        tx.BlockNumber,
			Timestamp:          tx.Timestamp.Unix(),
			GasUsed:            tx.GasUsed,
			GasPriceWei:        tx.GasPriceWei,
			TransactionFeeEth:  float64(tx.TransactionFeeEth.Float64),
			TransactionFeeUsdt: float64(tx.TransactionFeeUsdt.Float64),
			EthUsdtPrice:       float64(tx.EthUsdtPrice.Float64),
		}
		response = append(response, txResp)
	}

	ctx.JSON(http.StatusOK, response)
}

func (th *TransactionHandler) getTransactionByTimestamp(ctx *gin.Context) {
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	if startStr == "" || endStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Start and end timestamps are required"})
		return
	}

	// Parse timestamps as Unix time (seconds)
	startUnix, err := utils.ParseUnixTime(startStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp. Use Unix time in seconds."})
		return
	}
	endUnix, err := utils.ParseUnixTime(endStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp. Use Unix time in seconds."})
		return
	}

	// Convert Unix timestamps to time.Time
	startTime := time.Unix(startUnix, 0)
	endTime := time.Unix(endUnix, 0)

	if endTime.Before(startTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "End timestamp must be after start timestamp"})
		return
	}

	params := db.GetTransactionsByTimeRangeParams{
		Timestamp:   startTime,
		Timestamp_2: endTime,
	}
	transactions, err := th.txDbQuery.GetTransactionsByTimeRange(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var response []TransactionResponse
	for _, tx := range transactions {
		txResp := TransactionResponse{
			TransactionHash:    tx.TransactionHash,
			BlockNumber:        tx.BlockNumber,
			Timestamp:          tx.Timestamp.Unix(),
			GasUsed:            tx.GasUsed,
			GasPriceWei:        tx.GasPriceWei,
			TransactionFeeEth:  float64(tx.TransactionFeeEth.Float64),
			TransactionFeeUsdt: float64(tx.TransactionFeeUsdt.Float64),
			EthUsdtPrice:       float64(tx.EthUsdtPrice.Float64),
		}
		response = append(response, txResp)
	}

	ctx.JSON(http.StatusOK, response)
}
