package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/postgres"
)

type balanceStorage struct {
	baseStorage
}

func NewBalanceStorage(pg *postgres.Postgres) *balanceStorage {
	return &balanceStorage{
		baseStorage{pg},
	}
}

func (r *balanceStorage) Save(ctx context.Context, balance *entity.Balance) error {
	sql, args, err := r.db.Builder.
		Insert("balance").
		Columns("player_name, currency, balance").
		Values(balance.PlayerName, balance.Currency, balance.Balance).
		Suffix("ON CONFLICT (player_name, currency) DO UPDATE SET balance = balance + ?", balance.Balance).
		ToSql()
	if err != nil {
		return fmt.Errorf("BalanceStorage - Save - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("BalanceStorage - Save - r.Exec: %w", err)
	}
	return nil
}

func (r *balanceStorage) Get(ctx context.Context, playerName string, currency string) (*entity.Balance, error) {
	sql, args, err := r.db.Builder.
		Select("player_name, currency, balance").
		From("balance").
		Where(sq.Eq{"player_name": playerName, "currency": currency}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("BalanceStorage - Get - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("BalanceStorage - Get - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Balance{}
		err = rows.Scan(&e.PlayerName, &e.Currency, &e.Balance)
		if err != nil {
			return nil, fmt.Errorf("BalanceStorage - Get - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}
