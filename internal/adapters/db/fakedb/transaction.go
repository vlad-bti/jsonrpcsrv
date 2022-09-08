package fakedb

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type transactionStorage struct {
	db map[string]*entity.Transaction
}

func NewTransactionStorage() *transactionStorage {
	return &transactionStorage{
		db: make(map[string]*entity.Transaction),
	}
}

func (r *transactionStorage) Save(ctx context.Context, trx *entity.Transaction) error {
	r.db[trx.TransactionRef] = trx
	return nil
}

func (r *transactionStorage) Get(ctx context.Context, transactionRef string) (*entity.Transaction, error) {
	if order, ok := r.db[transactionRef]; ok {
		return order, nil
	}
	return nil, nil
}
