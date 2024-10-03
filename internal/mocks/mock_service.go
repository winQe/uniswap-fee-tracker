package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockBatchDataProcessor is a mock implementation of the BatchDataProcessor interface.
type MockBatchDataProcessor struct {
	mock.Mock
	// Channel to signal that ProcessBatchJob was called
	CalledChan chan struct{}
}

func NewMockBatchDataProcessor() *MockBatchDataProcessor {
	return &MockBatchDataProcessor{
		CalledChan: make(chan struct{}, 1),
	}
}

func (m *MockBatchDataProcessor) ProcessBatchJob(jobID string, startTime, endTime int64) error {
	m.Called(jobID, startTime, endTime)
	// Signal that the method was called
	m.CalledChan <- struct{}{}
	return nil
}
