package domains

import "eWallet/internal/shema"

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	GenerateWallet() (string, float64, error)
	Transaction(from string, to string, amount float64) error
	GetHistory(id string) ([]shema.HistoryTransfers, error)
	GetStatus(id string) (string, float64, error)
}
