package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	db "github.com/winQe/uniswap-fee-tracker/internal/db/sqlc"
)

// MockQuerier is a mock implementation of the db.Querier interface.
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) GetTransactionByHash(ctx context.Context, hash string) (db.Transactions, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(db.Transactions), args.Error(1)
}

func (m *MockQuerier) GetLatestTransactions(ctx context.Context, limit int32) ([]db.Transactions, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Transactions), args.Error(1)
}

func (m *MockQuerier) GetTransactionsByTimeRange(ctx context.Context, params db.GetTransactionsByTimeRangeParams) ([]db.Transactions, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]db.Transactions), args.Error(1)
}

func (m *MockQuerier) GetTransactionsByBlockNumber(ctx context.Context, blockNumber int64) ([]db.Transactions, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).([]db.Transactions), args.Error(1)
}

func (m *MockQuerier) InsertTransaction(ctx context.Context, arg db.InsertTransactionParams) error {
	return nil
}
