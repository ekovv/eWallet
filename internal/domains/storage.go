package domains

import (
	"context"
	"eWallet/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	SaveWallet(ctx context.Context, id string, balance float64) error
	Close() error
	TakeWallet(ctx context.Context, id string) (string, float64, error)
	UpdateWallet(ctx context.Context, id string, balance float64) error
	SaveInfo(ctx context.Context, from string, to string, amount float64, time string) error
	GetInfo(ctx context.Context, id string) ([]shema.HistoryTransfers, error)
}
