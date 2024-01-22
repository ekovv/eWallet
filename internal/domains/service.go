package domains

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	GenerateWallet() (string, float64, error)
}
