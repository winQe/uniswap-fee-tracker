package cache

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// BatchJob represents the response structure for a batch job.
type BatchJob struct {
	ID        string `json:"id"`     // UUID of the batch job
	Status    string `json:"status"` // Status: pending, running, completed, failed
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Result    string `json:"result"`
}

// JobsCache keeps track of all pending batch jobs within the TTL
type JobsCache struct {
	*RedisCache
	keyPrefix  string
	expiryTime time.Duration
}

// ErrJobNotFound is returned when a job is not found in the cache.
var ErrJobNotFound = errors.New("job not found")

const jobDB = 2

func NewJobCache(addr, password string) *JobsCache {
	return &JobsCache{
		RedisCache: NewRedisCache(addr, password, 2),
		keyPrefix:  "batch_job",
		expiryTime: 30 * time.Minute,
	}
}

// SetJob stores a batch job in Redis with the given job ID.
func (jb *JobsCache) SetJob(jobID string, jobData []byte) error {
	key := "batch_job:" + jobID
	return jb.client.Set(jb.ctx, key, jobData, jb.expiryTime).Err()
}

// GetJob retrieves a batch job from Redis by job ID.
func (jb *JobsCache) GetJob(jobID string) ([]byte, error) {
	key := "batch_job:" + jobID
	jobData, err := jb.client.Get(jb.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return jobData, nil
}

// GetAllJobs retrieves all batch jobs from Redis.
func (jb *JobsCache) GetAllJobs() ([][]byte, error) {
	var jobs [][]byte
	iter := jb.client.Scan(jb.ctx, 0, "batch_job:*", 0).Iterator()
	for iter.Next(jb.ctx) {
		jobData, err := jb.client.Get(jb.ctx, iter.Val()).Bytes()
		if err != nil {
			continue // Skip if unable to get job data
		}
		jobs = append(jobs, jobData)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}
