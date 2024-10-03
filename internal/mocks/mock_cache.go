package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRateCache is a mock implementation of the RateCache interface.
type MockRateCache struct {
	mock.Mock
}

// GetRate mocks the GetRate method of RateCache.
func (m *MockRateCache) GetRate(timestamp time.Time) (float64, error) {
	args := m.Called(timestamp)
	return args.Get(0).(float64), args.Error(1)
}

// StoreRate mocks the StoreRate method of RateCache.
func (m *MockRateCache) StoreRate(timestamp time.Time, value float64) error {
	args := m.Called(timestamp, value)
	return args.Error(0)
}

// MockJobsStore is a mock implementation of the JobsStore interface.
type MockJobsStore struct {
	mock.Mock
}

func (m *MockJobsStore) SetJob(id string, data []byte) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockJobsStore) GetJob(id string) ([]byte, error) {
	args := m.Called(id)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockJobsStore) GetAllJobs() ([][]byte, error) {
	args := m.Called()
	return args.Get(0).([][]byte), args.Error(1)
}
