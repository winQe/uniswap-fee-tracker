package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

// TransactionResponse represents the JSON structure of a transaction in the API response.
// swagger:model
type TransactionResponse struct {
	// The hash of the transaction
	TransactionHash string `json:"transaction_hash"`
	// The block number where the transaction was included
	BlockNumber int64 `json:"block_number"`
	// The timestamp of the transaction (Unix epoch time in seconds)
	Timestamp int64 `json:"timestamp"`
	// The amount of gas used by the transaction
	GasUsed int64 `json:"gas_used"`
	// The gas price in Wei
	GasPriceWei int64 `json:"gas_price_wei"`
	// The transaction fee in Ether
	TransactionFeeEth float64 `json:"transaction_fee_eth"`
	// The transaction fee in USDT
	TransactionFeeUsdt float64 `json:"transaction_fee_usdt"`
	// The Ether to USDT price at the time of the transaction
	EthUsdtPrice float64 `json:"eth_usdt_price"`
}

// TransactionHandler handles transaction related CRUD logic
type TransactionHandler struct {
	txDbQuery db.Querier
}

// NewTransactionHandler initializes a new TransactionHandler with the given dependencies.
func NewTransactionHandler(txDbQuery db.Querier) *TransactionHandler {
	return &TransactionHandler{
		txDbQuery: txDbQuery,
	}
}

// getTransactionHash godoc
// @Summary Get transaction by hash
// @Description Retrieve a specific transaction using its hash.
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param hash path string true "Transaction Hash"
// @Success 200 {object} TransactionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/{hash} [get]
func (th *TransactionHandler) getTransactionHash(ctx *gin.Context) {
	txHash := ctx.Param("hash")
	txHash = utils.SanitizeTransactionHash(txHash)

	if txHash == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid or missing transaction hash"})
		return
	}

	// Fetch transaction from the database
	transaction, err := th.txDbQuery.GetTransactionByHash(ctx, txHash)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Transaction not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
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

// getLatestTransactions godoc
// @Summary Get latest transactions
// @Description Retrieve the latest transactions with an optional limit.
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of transactions to retrieve" default(10)
// @Success 200 {array} TransactionResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/latest [get]
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
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		log.Printf("error getting latest txs %v", err)
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

// getTransactionByTimestamp godoc
// @Summary Get transactions within a timestamp range
// @Description Retrieve a list of transactions that occurred between the specified start and end Unix epoch timestamps.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param start query string true "Start timestamp in Unix epoch seconds"
// @Param end query string true "End timestamp in Unix epoch seconds"
// @Success 200 {array} TransactionResponse "List of transactions"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /transactions [get]
func (th *TransactionHandler) getTransactionByTimestamp(ctx *gin.Context) {
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	if startStr == "" || endStr == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Start and end timestamps are required"})
		return
	}

	// Parse timestamps as Unix time (seconds)
	startUnix, err := utils.ParseUnixTime(startStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid start timestamp. Use Unix time in seconds."})
		return
	}
	endUnix, err := utils.ParseUnixTime(endStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid end timestamp. Use Unix time in seconds."})
		return
	}

	// Convert Unix timestamps to time.Time
	startTime := time.Unix(startUnix, 0)
	endTime := time.Unix(endUnix, 0)

	if endTime.Before(startTime) {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "End timestamp must be after start timestamp"})
		return
	}

	params := db.GetTransactionsByTimeRangeParams{
		Timestamp:   startTime,
		Timestamp_2: endTime,
	}
	transactions, err := th.txDbQuery.GetTransactionsByTimeRange(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
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
