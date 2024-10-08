// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	GetLatestTransactions(ctx context.Context, limit int32) ([]Transactions, error)
	GetTransactionByHash(ctx context.Context, transactionHash string) (Transactions, error)
	GetTransactionsByBlockNumber(ctx context.Context, blockNumber int64) ([]Transactions, error)
	GetTransactionsByTimeRange(ctx context.Context, arg GetTransactionsByTimeRangeParams) ([]Transactions, error)
	InsertTransaction(ctx context.Context, arg InsertTransactionParams) error
}

var _ Querier = (*Queries)(nil)
