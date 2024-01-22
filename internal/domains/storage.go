package domains

type Storage interface {
	SaveWallet(id string, balance float64) error
}
