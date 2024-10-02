package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBAddress       string `mapstructure:"DB_ADDRESS"`
	DBName          string `mapstructure:"DB_NAME"`
	RedisURL        string `mapstructure:"REDIS_URL"`
	EtherscanAPIKey string `mapstructure:"ETHERSCAN_API_KEY"`
	ServerPort      string `mapstructure:"SERVER_PORT"`
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
