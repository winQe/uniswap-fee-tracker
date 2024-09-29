package cache

import (
	"fmt"
	"strconv"
	"time"
)

type RateCache struct {
	*RedisCache
}

const rateDB = 0

func NewRateCache(addr, password string) *RateCache {
	return &RateCache{
		NewRedisCache(addr, password, rateDB),
	}
}

// StorePrice stores a price value for the given key (e.g., "eth-usdt") which has TTL of 1 minute
func (rc *RateCache) StoreRate(key string, price float64) error {
	return rc.client.Set(rc.ctx, key, price, time.Minute).Err()
}

// GetPrice retrieves the rate value for the given key
func (rc *RedisCache) GetRate(key string) (float64, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error could not retrieve key %s: %w", key, err)
	}

	price, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("error could not parse value for key %s: %w", key, err)
	}

	return price, nil
}
