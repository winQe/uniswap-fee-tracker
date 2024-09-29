// Package client provides utilities to interact with external APIs such as Etherscan, Binance, etc
package client

import "time"

// PriceClient defines the interface for fetching price data. Mostly for dependency injection
type PriceClient interface {
	GetETHUSDT(timestamp time.Time) (*KlineData, error)
}
