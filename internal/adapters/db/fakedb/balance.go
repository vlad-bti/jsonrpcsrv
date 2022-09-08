package fakedb

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type balanceStorage struct {
	db map[recordType]*entity.Balance
}

type recordType struct {
	PlayerName string
	Currency   string
}

func NewBalanceStorage() *balanceStorage {
	return &balanceStorage{
		db: make(map[recordType]*entity.Balance),
	}
}

func (r *balanceStorage) Save(ctx context.Context, balance *entity.Balance) error {
	r.db[recordType{PlayerName: balance.PlayerName, Currency: balance.Currency}] = balance
	return nil
}

func (r *balanceStorage) Get(ctx context.Context, playerName string, currency string) (*entity.Balance, error) {
	if record, ok := r.db[recordType{PlayerName: playerName, Currency: currency}]; ok {
		return record, nil
	}
	return nil, nil
}
