package dto

type GetBalanceParams struct {
	CallerId             int    `json:"callerId" binding:"required"`
	PlayerName           string `json:"playerName" binding:"required"`
	Currency             string `json:"currency" binding:"required"`
	GameId               string `json:"gameId"`
	SessionId            string `json:"sessionId"`
	SessionAlternativeId string `json:"sessionAlternativeId"`
	BonusId              string `json:"bonusId"`
}

type SpinDetailsType struct {
	BetType string `json:"betType"`
	WinType string `json:"winType"`
}

type WithdrawAndDepositParams struct {
	CallerId             int             `json:"callerId" binding:"required"`
	PlayerName           string          `json:"playerName" binding:"required"`
	Withdraw             int             `json:"withdraw" binding:"required"`
	Deposit              int             `json:"deposit" binding:"required"`
	Currency             string          `json:"currency" binding:"required"`
	TransactionRef       string          `json:"transactionRef" binding:"required"`
	GameRoundRef         string          `json:"gameRoundRef"`
	GameId               string          `json:"gameId"`
	Source               string          `json:"source"`
	Reason               string          `json:"reason"`
	SessionId            string          `json:"sessionId"`
	SessionAlternativeId string          `json:"sessionAlternativeId"`
	SpinDetails          SpinDetailsType `json:"spinDetails"`
	BonusId              string          `json:"bonusId"`
	ChargeFreerounds     int             `json:"chargeFreerounds"`
}

type RollbackTransactionParams struct {
	CallerId             int    `json:"callerId" binding:"required"`
	PlayerName           string `json:"playerName" binding:"required"`
	TransactionRef       string `json:"transactionRef" binding:"required"`
	GameId               string `json:"gameId"`
	SessionId            string `json:"sessionId"`
	SessionAlternativeId string `json:"sessionAlternativeId"`
	RoundId              string `json:"roundId"`
}
