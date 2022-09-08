package json_rpc

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type WithdrawAndDepositResponse struct {
	NewBalance     int    `json:"newBalanceValue"`
	TransactionId  string `json:"transactionId"`
	FreeroundsLeft int    `json:"freeroundsLeft,omitempty"`
}
