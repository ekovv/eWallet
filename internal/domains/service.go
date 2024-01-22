package domains

type Service interface {
	GenerateWallet() (string, error)
}
