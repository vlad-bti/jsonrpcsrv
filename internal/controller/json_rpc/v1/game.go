package json_rpc

import (
	"context"
	"encoding/json"

	"github.com/vlad-bti/jrpc2"
	"github.com/vlad-bti/jsonrpcsrv/internal/controller/json_rpc/dto"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/usecase"
)

type GameUsecase interface {
	GetBalance(ctx context.Context, dto usecase.GetBalanceDTO) (*entity.Balance, error)
	WithdrawAndDeposit(ctx context.Context, dto usecase.WithdrawAndDepositDTO) error
	RollbackTransaction(ctx context.Context, dto usecase.RollbackTransactionDTO) error
}

type gameHandler struct {
	gameUsecase GameUsecase
}

func NewGameHandler(gameUsecase GameUsecase) *gameHandler {
	return &gameHandler{gameUsecase: gameUsecase}
}

func (g *gameHandler) Register(s *jrpc2.Server) {
	s.RegisterWithContext("getBalance", jrpc2.MethodWithContext{Method: g.GetBalance})
	s.RegisterWithContext("withdrawAndDeposit", jrpc2.MethodWithContext{Method: g.WithdrawAndDeposit})
	s.RegisterWithContext("rollbackTransaction", jrpc2.MethodWithContext{Method: g.RollbackTransaction})
}

func (g *gameHandler) GetBalance(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(dto.GetBalanceParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	usecaseDTO := usecase.GetBalanceDTO{
		PlayerName: p.PlayerName,
		Currency:   p.Currency,
	}
	resp, err := g.gameUsecase.GetBalance(ctx, usecaseDTO)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrUnknown,
			Message: jrpc2.ErrorMsg(err.Error()),
		}
	}
	var balance int
	if resp != nil {
		balance = resp.Balance
	}
	return BalanceResponse{Balance: balance}, nil
}

func (g *gameHandler) WithdrawAndDeposit(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(dto.WithdrawAndDepositParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Deposit < 0 {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrNegativeDepositCode,
			Message: ErrNegativeDepositCodeName,
		}
	}
	if p.Withdraw < 0 {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrNegativeWithdrawalCode,
			Message: ErrNegativeWithdrawalCodeName,
		}
	}
	if !entity.IsValidCurrency(p.Currency) {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrIllegalCurrencyCode,
			Message: ErrIllegalCurrencyCodeName,
		}
	}
	usecaseDTO := usecase.WithdrawAndDepositDTO{
		PlayerName:       p.PlayerName,
		Withdraw:         p.Withdraw,
		Deposit:          p.Deposit,
		Currency:         p.Currency,
		TransactionRef:   p.TransactionRef,
		ChargeFreerounds: p.ChargeFreerounds,
	}
	err := g.gameUsecase.WithdrawAndDeposit(ctx, usecaseDTO)
	if err != nil {
		if err == entity.ErrNotEnoughMoney {
			return nil, &jrpc2.ErrorObject{
				Code:    ErrNotEnoughMoneyCode,
				Message: ErrNotEnoughMoneyCodeName,
			}
		}
		return nil, &jrpc2.ErrorObject{
			Code:    ErrUnknown,
			Message: jrpc2.ErrorMsg(err.Error()),
		}
	}

	balanceDTO := usecase.GetBalanceDTO{
		PlayerName: p.PlayerName,
		Currency:   p.Currency,
	}
	resp, err := g.gameUsecase.GetBalance(ctx, balanceDTO)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrUnknown,
			Message: jrpc2.ErrorMsg(err.Error()),
		}
	}
	var balance int
	if resp != nil {
		balance = resp.Balance
	}
	return WithdrawAndDepositResponse{NewBalance: balance, TransactionId: p.TransactionRef}, nil
}

func (g *gameHandler) RollbackTransaction(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(dto.RollbackTransactionParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	usecaseDTO := usecase.RollbackTransactionDTO{
		TransactionRef: p.TransactionRef,
	}
	err := g.gameUsecase.RollbackTransaction(ctx, usecaseDTO)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    ErrUnknown,
			Message: jrpc2.ErrorMsg(err.Error()),
		}
	}
	return nil, nil
}
