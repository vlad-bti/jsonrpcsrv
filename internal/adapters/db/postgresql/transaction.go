package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/postgres"
)

type transactionStorage struct {
	baseStorage
}

func NewTransactionStorage(pg *postgres.Postgres) *transactionStorage {
	return &transactionStorage{
		baseStorage{pg},
	}
}

func (r *transactionStorage) Save(ctx context.Context, trx *entity.Transaction) error {
	sql, args, err := r.db.Builder.
		Insert("transaction").
		Columns("player_name, withdraw, deposit, currency, transaction_ref, charge_free_rounds, status").
		Values(
			trx.PlayerName,
			trx.Withdraw,
			trx.Deposit,
			trx.Currency,
			trx.TransactionRef,
			trx.ChargeFreerounds,
			trx.Status).
		Suffix("ON CONFLICT (transaction_ref) DO UPDATE SET player_name = ?, withdraw = ?, deposit = ?, currency = ?, charge_free_rounds = ?, status = ?",
			trx.PlayerName,
			trx.Withdraw,
			trx.Deposit,
			trx.Currency,
			trx.ChargeFreerounds,
			trx.Status).
		ToSql()
	if err != nil {
		return fmt.Errorf("TransactionStorage - Save - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TransactionStorage - Save - r.Exec: %w", err)
	}
	return nil
}

func (r *transactionStorage) Get(ctx context.Context, transactionRef string) (*entity.Transaction, error) {
	sql, args, err := r.db.Builder.
		Select("player_name, withdraw, deposit, currency, transaction_ref, charge_free_rounds, status").
		From("transaction").
		Where(sq.Eq{"transaction_ref": transactionRef}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionStorage - Get - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionStorage - Get - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Transaction{}
		err = rows.Scan(
			&e.PlayerName,
			&e.Withdraw,
			&e.Deposit,
			&e.Currency,
			&e.TransactionRef,
			&e.ChargeFreerounds,
			&e.Status)
		if err != nil {
			return nil, fmt.Errorf("TransactionStorage - Get - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}
