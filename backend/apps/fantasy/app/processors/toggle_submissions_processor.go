package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"time"
)

type ToggleSubmissionsProcessorInterface interface {
	ToggleSubmissions(requester *models.Player, enabled bool) (*models.GameView, error)
}

type ToggleSubmissionsProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ToggleSubmissionsProcessor) New(gameRepository repositories.GameRepository) *ToggleSubmissionsProcessor {
	return &ToggleSubmissionsProcessor{
		gameRepository,
	}
}

func (P *ToggleSubmissionsProcessor) ToggleSubmissions(requester *models.Player, enabled bool) (*models.GameView, error) {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	admin := game.GetPlayerByUUID(requester.UUID)
	if admin == nil || admin.Class != models.PLAYER_CLASS_ADMIN {
		return nil, ErrForbidden
	}

	game.SubmissionsEnabled = enabled
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	err = messages.SubmissionsEnabledMessage(game.Key, enabled)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
