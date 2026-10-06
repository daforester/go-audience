package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"time"
)

type SyncVolumeProcessorInterface interface {
	SyncVolume(requester *models.Player, volume int) (*models.GameView, error)
}

type SyncVolumeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SyncVolumeProcessor) New(gameRepository repositories.GameRepository) *SyncVolumeProcessor {
	return &SyncVolumeProcessor{
		gameRepository,
	}
}

func (P *SyncVolumeProcessor) SyncVolume(requester *models.Player, volume int) (*models.GameView, error) {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(requester.UUID)
	if player == nil || (player.Class != models.PLAYER_CLASS_ADMIN && player.Class != models.PLAYER_CLASS_DISPLAY) {
		return nil, ErrForbidden
	}

	game.Volume = min(max(volume, 0), 100)
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game.ViewFor(player), err
	}

	err = messages.SetVolumeMessage(game.Key, game.Volume)
	if err != nil {
		return game.ViewFor(player), err
	}

	return game.ViewFor(player), nil
}
