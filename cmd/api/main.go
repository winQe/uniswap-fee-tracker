package main

import (
	"log"

	"github.com/winQe/uniswap-fee-tracker/internal/server"
	"github.com/winQe/uniswap-fee-tracker/internal/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	server := server.NewServer(config.ServerPort)
	server.Run()
}
