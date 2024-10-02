package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUser              string `mapstructure:"DB_USER"`
	DBPassword          string `mapstructure:"DB_PASSWORD"`
	DBAddress           string `mapstructure:"DB_ADDRESS"`
	DBPort              string `mapstructure:"DB_PORT"`
	DBName              string `mapstructure:"DB_NAME"`
	RedisURL            string `mapstructure:"REDIS_URL"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	EtherscanAPIKey     string `mapstructure:"ETHERSCAN_API_KEY"`
	ServerPort          string `mapstructure:"SERVER_PORT"`
	WETHUSDCPoolAddress string `mapstructure:"WETH_USDT_POOL_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".") // can call multiple times to add more search path
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
