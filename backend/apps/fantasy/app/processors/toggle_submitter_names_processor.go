package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"time"
)

type ToggleSubmitterNamesProcessorInterface interface {
	ToggleSubmitterNames(requester *models.Player, enabled bool) (*models.GameView, error)
}

type ToggleSubmitterNamesProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ToggleSubmitterNamesProcessor) New(gameRepository repositories.GameRepository) *ToggleSubmitterNamesProcessor {
	return &ToggleSubmitterNamesProcessor{
		gameRepository,
	}
}

// ToggleSubmitterNames changes whether the display, directors and audience see who wrote
// each submission. Everyone gets a fresh view because names appear across the game.
func (P *ToggleSubmitterNamesProcessor) ToggleSubmitterNames(requester *models.Player, enabled bool) (*models.GameView, error) {
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

	game.ShowSubmitterNames = enabled
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	err = messages.GameUpdateMessage(game)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
