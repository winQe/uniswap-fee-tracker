package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisCache abstracts Redis client
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis client instance. Can be called for different DBs
func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("error connecting to Redis: %v", err)
	}

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Close gracefully closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
