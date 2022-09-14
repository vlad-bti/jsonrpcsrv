package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/postgres"
)

type playerStorage struct {
	baseStorage
}

func NewPlayerStorage(pg *postgres.Postgres) *playerStorage {
	return &playerStorage{
		baseStorage{pg},
	}
}

func (r *playerStorage) Save(ctx context.Context, player *entity.Player) error {
	sql, args, err := r.db.Builder.
		Insert("player").
		Columns("player_name, free_rounds").
		Values(player.PlayerName, player.Freerounds).
		Suffix("ON CONFLICT (player_name) DO UPDATE SET free_rounds = ?", player.Freerounds).
		ToSql()
	if err != nil {
		return fmt.Errorf("PlayerStorage - Save - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PlayerStorage - Save - r.Exec: %w", err)
	}
	return nil
}

func (r *playerStorage) Get(ctx context.Context, playerName string) (*entity.Player, error) {
	sql, args, err := r.db.Builder.
		Select("player_name, free_rounds").
		From("player").
		Where(sq.Eq{"player_name": playerName}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PlayerStorage - Get - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PlayerStorage - Get - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Player{}
		err = rows.Scan(&e.PlayerName, &e.Freerounds)
		if err != nil {
			return nil, fmt.Errorf("PlayerStorage - Get - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}
