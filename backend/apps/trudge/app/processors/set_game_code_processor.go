package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"errors"
	"strconv"
	"time"
)

type SetGameCodeProcessorInterface interface {
	SetGameCode(roomKey string, mode int) (*models.Game, error)
}

type SetGameCodeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SetGameCodeProcessor) New(gameRepository repositories.GameRepository) *SetGameCodeProcessor {
	return &SetGameCodeProcessor{
		gameRepository,
	}
}

func (P *SetGameCodeProcessor) SetGameCode(roomKey string, mode int) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.GameCode = mode
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.setGameCodeMessage(game.Key, mode)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *SetGameCodeProcessor) setGameCodeMessage(roomKey string, mode int) error {
	wsCollection := sockets.NewCollection(roomKey)
	wsCollection.SendText("CODE "+strconv.Itoa(mode), wsCollection.All()...)

	return nil
}
