package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
	"github.com/winQe/uniswap-fee-tracker/internal/domain"
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

	// Initialize LiveDataRecorder
	liveDataRecorder := service.NewLiveDataRecorder(dbQuerier, txManager)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal: %v. Initiating shutdown.", sig)
		cancel()
	}()

	// Run the LiveDataRecorder in a separate goroutine
	go liveDataRecorder.Run(ctx)

	// Keep the main function running until context is canceled
	<-ctx.Done()
	log.Println("Application has shut down gracefully.")
}
