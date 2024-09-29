package domain

import (
	"fmt"
	"log"
	"time"

	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
)

type PriceManager struct {
	rateCache   cache.RateStore
	priceClient client.PriceClient
}

// NewPriceManager creates a PriceManager for handling logic with getting ETH-USDT conversion rate
func NewPriceManager(rateStore cache.RateStore, priceClient client.PriceClient) *PriceManager {
	return &PriceManager{
		rateCache:   rateStore,
		priceClient: priceClient,
	}
}

// GetETHUSDTPrice retrieves the price of ETH to USDT.
// It first checks the cache, and if not found, it fetches from Binance SPOT API, then caches the result.
func (p *PriceManager) GetETHUSDT(timestamp time.Time) (float64, error) {
	// Check the cache first
	price, err := p.rateCache.GetRate("eth-usdt")
	// Cache hit
	if err == nil {
		return price, nil
	}

	// Cache miss, get price from external API
	klineData, err := p.priceClient.GetETHUSDT(timestamp)
	if err != nil {
		return 0, fmt.Errorf("could not get ETH to USDT price from external API: %w", err)
	}

	// Store in cache for future use
	err = p.rateCache.StoreRate("eth-usdt", klineData.ClosePrice)
	if err != nil {
		log.Printf("Warning: could not store price in cache: %v\n", err)
	}

	return klineData.ClosePrice, nil
}
