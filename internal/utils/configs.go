package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser              string
	DBPassword          string
	DBAddress           string
	DBPort              string
	DBName              string
	RedisURL            string
	RedisPassword       string
	EtherscanAPIKey     string
	ServerPort          string
	WETHUSDCPoolAddress string
}

// LoadConfig reads configuration from a .env file and environment variables.
func LoadConfig() (Config, error) {
	var config Config

	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found. Proceeding with environment variables.")
	}

	// Populate Config struct from environment variables
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBAddress = os.Getenv("DB_ADDRESS")
	config.DBPort = os.Getenv("DB_PORT")
	config.DBName = os.Getenv("DB_NAME")
	config.RedisURL = os.Getenv("REDIS_URL")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.EtherscanAPIKey = os.Getenv("ETHERSCAN_API_KEY")
	config.ServerPort = os.Getenv("SERVER_PORT")
	config.WETHUSDCPoolAddress = os.Getenv("WETH_USDT_POOL_ADDRESS")

	// Validate required fields
	if config.DBUser == "" {
		return config, fmt.Errorf("DB_USER is required")
	}
	if config.DBPassword == "" {
		return config, fmt.Errorf("DB_PASSWORD is required")
	}
	if config.DBAddress == "" {
		return config, fmt.Errorf("DB_ADDRESS is required")
	}
	if config.DBPort == "" {
		return config, fmt.Errorf("DB_PORT is required")
	}
	if config.DBName == "" {
		return config, fmt.Errorf("DB_NAME is required")
	}
	if config.RedisURL == "" {
		return config, fmt.Errorf("REDIS_URL is required")
	}
	if config.EtherscanAPIKey == "" {
		return config, fmt.Errorf("ETHERSCAN_API_KEY is required")
	}
	if config.ServerPort == "" {
		return config, fmt.Errorf("SERVER_PORT is required")
	}
	if config.WETHUSDCPoolAddress == "" {
		return config, fmt.Errorf("WETH_USDT_POOL_ADDRESS is required")
	}

	return config, nil
}
