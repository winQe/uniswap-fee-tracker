package cache

import "time"

// RateStore defines the interface for interacting with the rate cache.
// It allows storing and retrieving rate values based on timestamps.
type RateStore interface {
	StoreRate(timestamp time.Time, price float64) error
	GetRate(timestamp time.Time) (float64, error)
}

type JobsStore interface {
	SetJob(jobID string, jobData []byte) error
	GetJob(jobID string) ([]byte, error)
	GetAllJobs() ([][]byte, error)
}
