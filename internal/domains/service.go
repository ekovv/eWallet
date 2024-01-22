package domains

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	GenerateWallet() (string, float64, error)
	Transaction(from string, to string, amount float64) error
}
