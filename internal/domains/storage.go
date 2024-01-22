package domains

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	SaveWallet(id string, balance float64) error
	Close() error
}
