package domains

import (
	"context"
	"eWallet/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	GenerateWallet(ctx context.Context) (string, float64, error)
	Transaction(from string, to string, amount float64, ctx context.Context) error
	GetHistory(id string, ctx context.Context) ([]shema.HistoryTransfers, error)
	GetStatus(id string, ctx context.Context) (string, float64, error)
}
