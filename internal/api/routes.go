package api

import (
	"github.com/gin-gonic/gin"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
)

func RegisterRoutes(rg *gin.RouterGroup, priceManager domain.PriceManagerInterface, transactionHandler *TransactionHandler) {
	rg.GET("/transactions/:hash", transactionHandler.getTransactionHash)
	rg.GET("/transactions/latest", transactionHandler.getLatestTransactions)
	rg.GET("/transactions/", transactionHandler.getTransactionByTimestamp)
}
