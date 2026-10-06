package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"time"
)

type KickPlayerProcessorInterface interface {
	KickPlayer(admin *models.Player, uuid string) error
}

type KickPlayerProcessor struct {
	gameRepository repositories.GameRepository
}

func (P KickPlayerProcessor) New(gameRepository repositories.GameRepository) *KickPlayerProcessor {
	return &KickPlayerProcessor{
		gameRepository,
	}
}

func (P *KickPlayerProcessor) KickPlayer(requester *models.Player, uuid string) error {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	admin := game.GetPlayerByUUID(requester.UUID)
	if admin == nil || admin.Class != models.PLAYER_CLASS_ADMIN {
		return ErrForbidden
	}

	player := game.GetPlayerByUUID(uuid)
	if player == nil {
		return ErrPlayerNotFound
	}

	messages.KickPlayerMessage(player)

	err := messages.PlayerQuitMessage(player)

	if err != nil {
		return err
	}

	game.Updated = time.Now()
	game.RemovePlayer(player)
	err = P.gameRepository.Save(game)

	if err != nil {
		return err
	}

	return nil
}
