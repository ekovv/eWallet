package shema

type Wallet struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type Transfer struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type HistoryTransfers struct {
	Time   string  `json:"time"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
