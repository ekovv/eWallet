package handler

type Wallet struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type Transfer struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
