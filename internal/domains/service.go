package domains

import (
	"context"
	"eWallet/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	GenerateWallet(ctx context.Context) (string, float64, error)
	Transaction(ctx context.Context, from string, to string, amount float64) error
	GetHistory(ctx context.Context, id string) ([]shema.HistoryTransfers, error)
	GetStatus(ctx context.Context, id string) (string, float64, error)
}
