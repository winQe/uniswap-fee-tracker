package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateCache implements the RateStore interface.
type RateCache struct {
	*RedisCache
	sortedSetKey string
	ttl          time.Duration
}

const rateDB = 0

// NewRateCache creates a new RateCache instance.
func NewRateCache(addr, password string) RateStore {
	return &RateCache{
		RedisCache:   NewRedisCache(addr, password, rateDB),
		sortedSetKey: "rate_cache",    // key for redis sorted set
		ttl:          5 * time.Minute, // price expires after 5 minutes
	}
}

// StoreRate stores the current price with the given timestamp.
func (rc *RateCache) StoreRate(timestamp time.Time, price float64) error {
	ts := timestamp.Unix()
	_, err := rc.client.ZAdd(rc.ctx, rc.sortedSetKey, redis.Z{
		Score:  float64(ts),
		Member: price,
	}).Result()
	if err != nil {
		return fmt.Errorf("error adding rate to sorted set: %w", err)
	}

	// Set the TTL for the sorted set key
	// Reset the TTL every time a new rate is stored to keep the key alive as long as data is being added
	err = rc.client.Expire(rc.ctx, rc.sortedSetKey, rc.ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration on sorted set: %w", err)
	}

	return nil
}

// GetRate retrieves the price within a 5-minute range of the given timestamp.
func (rc *RateCache) GetRate(timestamp time.Time) (float64, error) {
	ts := timestamp.Unix()

	// Scores to find prices within one minute range of the provided timestamp
	minScore := float64(ts - 60*5)
	maxScore := float64(ts + 60*5)

	// Retrieve members with their scores within the specified score range
	zRange, err := rc.client.ZRangeByScoreWithScores(rc.ctx, rc.sortedSetKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", maxScore),
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("error querying sorted set with scores: %w", err)
	}

	if len(zRange) == 0 {
		return 0, fmt.Errorf("no rate found within 5-minute range of timestamp %v", timestamp)
	}

	// Return the first rate found within the range
	firstMember := zRange[0].Member.(string)
	price, err := strconv.ParseFloat(firstMember, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing price: %w", err)
	}

	return price, nil
}
