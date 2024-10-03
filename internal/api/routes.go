package api

import (
	"github.com/gin-gonic/gin"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
)

func RegisterRoutes(rg *gin.RouterGroup, priceManager domain.PriceManagerInterface, transactionHandler *TransactionHandler, batchJobHandler *BatchJobHandler) {
	rg.GET("/transactions/:hash", transactionHandler.getTransactionHash)
	rg.GET("/transactions/latest", transactionHandler.getLatestTransactions)
	rg.GET("/transactions/", transactionHandler.getTransactionByTimestamp)

	rg.POST("/batch-jobs", batchJobHandler.CreateBatchJob)
	rg.GET("/batch-jobs/:id", batchJobHandler.GetBatchJob)
	rg.GET("/batch-jobs", batchJobHandler.ListBatchJobs)
}
