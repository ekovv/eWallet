package domains

import (
	"context"
	"eWallet/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	SaveWallet(id string, balance float64, ctx context.Context) error
	Close() error
	TakeWallet(id string, ctx context.Context) (string, float64, error)
	UpdateWallet(id string, balance float64, ctx context.Context) error
	SaveInfo(from string, to string, amount float64, time string, ctx context.Context) error
	GetInfo(id string, ctx context.Context) ([]shema.HistoryTransfers, error)
}
