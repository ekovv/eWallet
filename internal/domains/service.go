package domains

type Service interface {
	GenerateWallet() (string, float64, error)
}
