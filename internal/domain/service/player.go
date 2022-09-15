package service

import (
	"context"

	"github.com/vlad-bti/jsonrpcsrv/internal/domain/entity"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
)

type PlayerStorage interface {
	Save(ctx context.Context, player *entity.Player) error
	Get(ctx context.Context, playerName string) (*entity.Player, error)
}

type playerService struct {
	storage PlayerStorage
	log     *logger.Logger
}

func NewPlayerService(log *logger.Logger, storage PlayerStorage) *playerService {
	return &playerService{storage: storage, log: log}
}

func (s *playerService) GetPlayer(ctx context.Context, playerName string) (*entity.Player, error) {
	return s.storage.Get(ctx, playerName)
}

func (s *playerService) ChangeFreerounds(ctx context.Context, playerName string, value int) error {
	player, err := s.storage.Get(ctx, playerName)
	if err != nil {
		s.log.Error("PlayerService - ChangeFreerounds - s.storage.Get: %v; playerName=%v", err, playerName)
		return err
	}
	if player == nil && value < 0 || player != nil && player.Freerounds+value < 0 {
		return entity.ErrNotEnoughFreerounds
	}
	if player == nil {
		player = &entity.Player{
			PlayerName: playerName,
		}
	}
	player.Freerounds += value
	return s.storage.Save(ctx, player)
}
