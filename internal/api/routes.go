package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, transactionHandler *TransactionHandler, batchJobHandler *BatchJobHandler) {
	// Register transactions handlers
	rg.GET("/transactions/:hash", transactionHandler.getTransactionHash)
	rg.GET("/transactions/latest", transactionHandler.getLatestTransactions)
	rg.GET("/transactions/", transactionHandler.getTransactionByTimestamp)

	// Register batch jobs handler
	rg.POST("/batch-jobs", batchJobHandler.CreateBatchJob)
	rg.GET("/batch-jobs/:id", batchJobHandler.GetBatchJob)
	rg.GET("/batch-jobs", batchJobHandler.ListBatchJobs)
}
