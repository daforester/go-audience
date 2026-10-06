package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
	"time"
)

type SyncVolumeProcessorInterface interface {
	SyncVolume(roomKey string, volume int) (*models.Game, error)
}

type SyncVolumeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SyncVolumeProcessor) New(gameRepository repositories.GameRepository) *SyncVolumeProcessor {
	return &SyncVolumeProcessor{
		gameRepository,
	}
}

func (P *SyncVolumeProcessor) SyncVolume(roomKey string, volume int) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Volume = volume
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = messages.SetVolumeMessage(game.Key, volume)
	if err != nil {
		return game, err
	}

	return game, nil
}
