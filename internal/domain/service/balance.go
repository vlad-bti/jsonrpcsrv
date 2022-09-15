package service

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
)

type BalanceStorage interface {
	Save(ctx context.Context, balance *entity.Balance) error
	Get(ctx context.Context, playerName string, currency string) (*entity.Balance, error)
}

type balanceService struct {
	storage BalanceStorage
	log     *logger.Logger
}

func NewBalanceService(log *logger.Logger, storage BalanceStorage) *balanceService {
	return &balanceService{storage: storage, log: log}
}

func (s *balanceService) GetBalance(ctx context.Context, playerName string, currency string) (*entity.Balance, error) {
	return s.storage.Get(ctx, playerName, currency)
}

func (s *balanceService) Deposit(ctx context.Context, playerName string, currency string, value int) error {
	balance, err := s.storage.Get(ctx, playerName, currency)
	if err != nil {
		s.log.Error("BalanceService - Deposit - s.storage.Get: %v; playerName=%v, currency=%v",
			err,
			playerName,
			currency,
		)
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
		s.log.Error("BalanceService - Deposit - s.storage.Save: %v; playerName=%v, currency=%v, balance=%v",
			err,
			playerName,
			currency,
			balance.Balance,
		)
		return err
	}
	return nil
}

func (s *balanceService) Withdraw(ctx context.Context, playerName string, currency string, value int) error {
	balance, err := s.storage.Get(ctx, playerName, currency)
	if err != nil {
		s.log.Error("BalanceService - Withdraw - s.storage.Get: %v; playerName=%v, currency=%v",
			err,
			playerName,
			currency,
		)
		return err
	}
	if balance == nil || balance.Balance < value {
		return entity.ErrNotEnoughMoney
	}
	balance.Balance -= value
	return s.storage.Save(ctx, balance)
}
