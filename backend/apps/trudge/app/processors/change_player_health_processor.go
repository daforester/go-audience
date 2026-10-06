package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
	"time"
)

type ChangePlayerHealthProcessorInterface interface {
	ChangePlayerHealth(roomKey string, playerUUID string, adjustment int) (*models.Game, error)
}

type ChangePlayerHealthProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangePlayerHealthProcessor) New(gameRepository repositories.GameRepository) *ChangePlayerHealthProcessor {
	return &ChangePlayerHealthProcessor{
		gameRepository,
	}
}

func (P *ChangePlayerHealthProcessor) ChangePlayerHealth(roomKey string, playerUUID string, adjustment int) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	var player *models.Player
	var index int
	for i, p := range game.Players {
		if p.UUID == playerUUID {
			player = p
			index = i
		}
	}

	if player == nil {
		return nil, errors.New("player not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player.Health += adjustment
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = messages.HealthUpdateMessage(roomKey, player, uint(index), player.Health)
	if err != nil {
		return game, err
	}

	return game, nil
}
