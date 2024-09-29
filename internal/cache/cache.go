package cache

// RateStore defines the interface for interacting with rate cache. Mostly for dependency injection
type RateStore interface {
	StoreRate(key string, price float64) error
	GetRate(key string) (float64, error)
}
