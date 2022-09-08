package service

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type TransactionStorage interface {
	Save(ctx context.Context, trx *entity.Transaction) error
	Get(ctx context.Context, transactionRef string) (*entity.Transaction, error)
}

type FakeDB interface {
	Begin(ctx context.Context)
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
}

type transactionService struct {
	storage TransactionStorage
	fakeDB  FakeDB
}

func NewTransactionService(storage TransactionStorage, fakeDB FakeDB) *transactionService {
	return &transactionService{storage: storage, fakeDB: fakeDB}
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

func (s *transactionService) Begin(ctx context.Context) {
	s.fakeDB.Begin(ctx)
}

func (s *transactionService) Commit(ctx context.Context) {
	s.fakeDB.Commit(ctx)
}

func (s *transactionService) Rollback(ctx context.Context) {
	s.fakeDB.Rollback(ctx)
}
