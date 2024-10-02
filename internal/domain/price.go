package domain

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/winQe/uniswap-fee-tracker/internal/cache"
	"github.com/winQe/uniswap-fee-tracker/internal/client"
)

type conversionRate struct {
	timestamp time.Time
	rate      float64
}

// PriceManager aggregates the price from both cached rates and external APIs
type PriceManager struct {
	rateCache   cache.RateStore
	priceClient client.PriceClient
	lastPrice   conversionRate
	mu          sync.RWMutex // Protects access to lastPrice
}

// NewPriceManager creates a PriceManager for handling logic with getting ETH-USDT conversion rate
func NewPriceManager(rateStore cache.RateStore, priceClient client.PriceClient) *PriceManager {
	return &PriceManager{
		rateCache:   rateStore,
		priceClient: priceClient,
		lastPrice:   conversionRate{},
	}
}

// GetETHUSDTPrice retrieves the price of ETH to USDT.
// It first checks the lastPrice, then the cache, and finally fetches from the external API if needed.
func (p *PriceManager) GetETHUSDT(timestamp time.Time) (float64, error) {
	// Define a validity window for lastPrice (e.g., within the same minute)
	const validityDuration = time.Minute

	// Attempt to read the lastPrice with a read lock
	p.mu.RLock()
	if !p.lastPrice.timestamp.IsZero() && timestamp.Sub(p.lastPrice.timestamp) <= validityDuration {
		rate := p.lastPrice.rate
		p.mu.RUnlock()
		return rate, nil
	}
	p.mu.RUnlock()

	// Cache miss or lastPrice is stale, proceed to check the cache
	price, err := p.rateCache.GetRate(timestamp)
	if err == nil {
		// Update lastPrice with the fetched rate
		p.mu.Lock()
		p.lastPrice = conversionRate{
			timestamp: timestamp,
			rate:      price,
		}
		p.mu.Unlock()
		return price, nil
	}

	// Cache miss, fetch from external API
	klineData, err := p.priceClient.GetETHUSDT(timestamp)
	if err != nil {
		return 0, fmt.Errorf("could not get ETH to USDT price from external API: %w", err)
	}

	// Store the fetched rate in the cache
	err = p.rateCache.StoreRate(timestamp, klineData.ClosePrice)
	if err != nil {
		log.Printf("Warning: could not store price in cache: %v\n", err)
	}

	// Update lastPrice with the new rate
	p.mu.Lock()
	p.lastPrice = conversionRate{
		timestamp: timestamp,
		rate:      klineData.ClosePrice,
	}
	p.mu.Unlock()

	return klineData.ClosePrice, nil
}
