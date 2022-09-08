package json_rpc

import "github.com/vlad-bti/jrpc2"

const (
	_ jrpc2.ErrorCode = iota
	ErrNotEnoughMoneyCode
	ErrIllegalCurrencyCode
	ErrNegativeDepositCode
	ErrNegativeWithdrawalCode
	ErrSpendingBudgetExceeded

	ErrUnknown = 1000
)

const (
	ErrNotEnoughMoneyCodeName     = "ErrNotEnoughMoneyCode"
	ErrIllegalCurrencyCodeName    = "ErrIllegalCurrencyCode"
	ErrNegativeDepositCodeName    = "ErrNegativeDepositCode"
	ErrNegativeWithdrawalCodeName = "ErrNegativeWithdrawalCode"
	ErrSpendingBudgetExceededName = "ErrSpendingBudgetExceeded"
)
