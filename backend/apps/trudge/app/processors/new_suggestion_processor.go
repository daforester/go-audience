package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
)

type NewSuggestionProcessorInterface interface {
	NewSuggestion(roomKey string, suggestion string) (*models.Game, error)
}

type NewSuggestionProcessor struct {
	gameRepository repositories.GameRepository
}

func (P NewSuggestionProcessor) New(gameRepository repositories.GameRepository) *NewSuggestionProcessor {
	return &NewSuggestionProcessor{
		gameRepository,
	}
}

func (P *NewSuggestionProcessor) NewSuggestion(roomKey string, suggestion string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	err := messages.NewSuggestionMessage(game.Key, suggestion)
	if err != nil {
		return game, err
	}

	return game, nil
}
