package entity

type Transaction struct {
	PlayerName       string
	Withdraw         int
	Deposit          int
	Currency         string
	TransactionRef   string
	ChargeFreerounds int
	Status           int
}

const (
	TRANSACTION_DEFAULT int = iota
	TRANSACTION_REVERTED
)
