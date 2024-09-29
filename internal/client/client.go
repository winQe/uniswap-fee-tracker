// Package client provides utilities to interact with external APIs such as Etherscan, Binance, etc
package client

import "time"

// KlineData is the return type of PriceClient GetETHUSDT
type KlineData struct {
	ClosePrice float64
}

// PriceClient defines the interface for fetching price data. Mostly for dependency injection
type PriceClient interface {
	GetETHUSDT(timestamp time.Time) (*KlineData, error)
}
