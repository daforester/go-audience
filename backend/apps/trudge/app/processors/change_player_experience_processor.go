package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
	"time"
)

type ChangePlayerExperienceProcessorInterface interface {
	ChangePlayerExperience(roomKey string, playerUUID string, adjustment int) (*models.Game, error)
}

type ChangePlayerExperienceProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangePlayerExperienceProcessor) New(gameRepository repositories.GameRepository) *ChangePlayerExperienceProcessor {
	return &ChangePlayerExperienceProcessor{
		gameRepository,
	}
}

func (P *ChangePlayerExperienceProcessor) ChangePlayerExperience(roomKey string, playerUUID string, adjustment int) (*models.Game, error) {
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

	player.Experience += adjustment
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = messages.ExperienceUpdateMessage(roomKey, player, uint(index), player.Experience)
	if err != nil {
		return game, err
	}

	return game, nil
}
