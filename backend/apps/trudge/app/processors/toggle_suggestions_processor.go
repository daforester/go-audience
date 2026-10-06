package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
	"time"
)

type ToggleSuggestionsProcessorInterface interface {
	ToggleSuggestions(roomKey string, enabled bool) (*models.Game, error)
}

type ToggleSuggestionsProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ToggleSuggestionsProcessor) New(gameRepository repositories.GameRepository) *ToggleSuggestionsProcessor {
	return &ToggleSuggestionsProcessor{
		gameRepository,
	}
}

func (P *ToggleSuggestionsProcessor) ToggleSuggestions(roomKey string, enabled bool) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.SuggestionsEnabled = enabled
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = messages.SuggestionsEnabledMessage(game.Key, enabled)
	if err != nil {
		return game, err
	}

	return game, nil
}
