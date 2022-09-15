package service

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
)

type TransactionStorage interface {
	Save(ctx context.Context, trx *entity.Transaction) error
	Get(ctx context.Context, transactionRef string) (*entity.Transaction, error)
}

type transactionService struct {
	storage TransactionStorage
	log     *logger.Logger
}

func NewTransactionService(log *logger.Logger, storage TransactionStorage) *transactionService {
	return &transactionService{storage: storage, log: log}
}

func (s *transactionService) GetTransaction(ctx context.Context, transactionRef string) (*entity.Transaction, error) {
	return s.storage.Get(ctx, transactionRef)
}

func (s *transactionService) AddTransaction(ctx context.Context, trx *entity.Transaction) error {
	return s.storage.Save(ctx, trx)
}

func (s *transactionService) RevertTransaction(ctx context.Context, trx *entity.Transaction) error {
	transaction, err := s.storage.Get(ctx, trx.TransactionRef)
	if err != nil {
		s.log.Error("TransactionService - RevertTransaction - s.storage.sGet: %v; TransactionRef=%v", err, trx.TransactionRef)
		return err
	}
	if transaction == nil {
		trx.Status = entity.TRANSACTION_REVERTED
		return s.storage.Save(ctx, trx)
	}
	if transaction.Status == entity.TRANSACTION_DEFAULT {
		transaction.Status = entity.TRANSACTION_REVERTED
		return s.storage.Save(ctx, transaction)
	}
	return nil
}
