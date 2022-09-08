package usecase

type GetBalanceDTO struct {
	CallerId             int
	PlayerName           string
	Currency             string
	GameId               string
	SessionId            string
	SessionAlternativeId string
	BonusId              string
}

type SpinDetailsType struct {
	BetType string
	WinType string
}

type WithdrawAndDepositDTO struct {
	CallerId             int
	PlayerName           string
	Withdraw             int
	Deposit              int
	Currency             string
	TransactionRef       string
	GameRoundRef         string
	GameId               string
	Source               string
	Reason               string
	SessionId            string
	SessionAlternativeId string
	SpinDetails          SpinDetailsType
	BonusId              string
	ChargeFreerounds     int
}

type RollbackTransactionDTO struct {
	CallerId             int
	PlayerName           string
	TransactionRef       string
	GameId               string
	SessionId            string
	SessionAlternativeId string
	RoundId              string
}
