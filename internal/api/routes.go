package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/winQe/uniswap-fee-tracker/docs"
)

func RegisterRoutes(rg *gin.RouterGroup, transactionHandler *TransactionHandler, batchJobHandler *BatchJobHandler) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Register transactions handlers
	rg.GET("/transactions/:hash", transactionHandler.getTransactionHash)
	rg.GET("/transactions/latest", transactionHandler.getLatestTransactions)
	rg.GET("/transactions/", transactionHandler.getTransactionByTimestamp)

	// Register batch jobs handler
	rg.POST("/batch-jobs", batchJobHandler.CreateBatchJob)
	rg.GET("/batch-jobs/:id", batchJobHandler.GetBatchJob)
	rg.GET("/batch-jobs", batchJobHandler.ListBatchJobs)

	// Register Swagger route
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
