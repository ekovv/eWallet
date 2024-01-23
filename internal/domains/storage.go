package domains

import "eWallet/internal/shema"

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	SaveWallet(id string, balance float64) error
	Close() error
	TakeWallet(id string) (string, float64, error)
	UpdateWallet(id string, balance float64) error
	SaveInfo(from string, to string, amount float64, time string) error
	GetInfo(id string) ([]shema.HistoryTransfers, error)
}
