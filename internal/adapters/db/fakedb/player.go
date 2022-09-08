package fakedb

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
)

type playerStorage struct {
	db map[string]*entity.Player
}

func NewPlayerStorage() *playerStorage {
	return &playerStorage{
		db: make(map[string]*entity.Player),
	}
}

func (r *playerStorage) Save(ctx context.Context, player *entity.Player) error {
	r.db[player.PlayerName] = player
	return nil
}

func (r *playerStorage) Get(ctx context.Context, playerName string) (*entity.Player, error) {
	if record, ok := r.db[playerName]; ok {
		return record, nil
	}
	return nil, nil
}
