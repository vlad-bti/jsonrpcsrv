package service

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type BalanceStorage interface {
	Save(ctx context.Context, balance *entity.Balance) error
	Get(ctx context.Context, playerName string, currency string) (*entity.Balance, error)
}

type balanceService struct {
	storage BalanceStorage
}

func NewBalanceService(storage BalanceStorage) *balanceService {
	return &balanceService{storage: storage}
}

func (s *balanceService) GetBalance(ctx context.Context, playerName string, currency string) (*entity.Balance, error) {
	return s.storage.Get(ctx, playerName, currency)
}

func (s *balanceService) Deposit(ctx context.Context, playerName string, currency string, value int) error {
	balance, err := s.storage.Get(ctx, playerName, currency)
	if err != nil {
		return err
	}
	if balance == nil {
		balance = &entity.Balance{
			PlayerName: playerName,
			Currency:   currency,
			Balance:    value,
		}
	} else {
		balance.Balance += value
	}
	err = s.storage.Save(ctx, balance)
	if err != nil {
		return err
	}
	return nil
}

func (s *balanceService) Withdraw(ctx context.Context, playerName string, currency string, value int) error {
	balance, err := s.storage.Get(ctx, playerName, currency)
	if err != nil {
		return err
	}
	if balance == nil || balance.Balance < value {
		return entity.ErrNotEnoughMoney
	}
	balance.Balance -= value
	return s.storage.Save(ctx, balance)
}
