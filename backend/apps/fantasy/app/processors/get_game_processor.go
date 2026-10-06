package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
)

type GetGameProcessorInterface interface {
	GetGame(requester *models.Player) (*models.GameView, *models.Player, error)
}

type GetGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P GetGameProcessor) New(gameRepository repositories.GameRepository) *GetGameProcessor {
	return &GetGameProcessor{
		gameRepository,
	}
}

// GetGame returns the stored player and the view of the game that player is allowed to see
func (P *GetGameProcessor) GetGame(requester *models.Player) (*models.GameView, *models.Player, error) {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(requester.UUID)
	if player == nil {
		return nil, nil, ErrPlayerNotFound
	}

	return game.ViewFor(player), player.Copy(), nil
}
