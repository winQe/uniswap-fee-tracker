package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/winQe/uniswap-fee-tracker/internal/api"
	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
	"github.com/winQe/uniswap-fee-tracker/internal/server"
	"github.com/winQe/uniswap-fee-tracker/internal/service"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Initializes DB connection pool
	connPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBAddress, config.DBPort, config.DBName))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer connPool.Close()

	// Initialize dbQuerier from sqlc
	dbQuerier := db.New(connPool)

	// Initialize all price related dependencies
	priceCache := cache.NewRateCache(config.RedisURL, config.RedisPassword)
	binanceClient := client.NewKlineClient()
	priceManager := domain.NewPriceManager(priceCache, binanceClient)

	// Initialize all transactions related dependencies
	etherscanClient := client.NewEtherscanClient(config.EtherscanAPIKey, config.WETHUSDCPoolAddress)
	txManager := domain.NewTransactionManager(etherscanClient, priceManager)

	// Initialize batch job relatd dependencies
	jobsCache := cache.NewJobCache(config.RedisURL, config.RedisPassword)
	batchDataProcessor := service.NewBatchDataProcessor(dbQuerier, jobsCache, txManager)

	txHandler := api.NewTransactionHandler(dbQuerier)
	batchDataHandler := *api.NewBatchJobHandler(dbQuerier, jobsCache, txManager, batchDataProcessor)
	server := server.NewServer(config.ServerPort, txHandler, &batchDataHandler)

	server.Run()
}
