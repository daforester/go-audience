package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"errors"
	"strconv"
	"time"
)

type SetGameCodeProcessorInterface interface {
	SetGameCode(mode int) (*models.Game, error)
}

type SetGameCodeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SetGameCodeProcessor) New(gameRepository repositories.GameRepository) *SetGameCodeProcessor {
	return &SetGameCodeProcessor{
		gameRepository,
	}
}

func (P *SetGameCodeProcessor) SetGameCode(mode int) (*models.Game, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()
	game.GameCode = mode
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.setGameCodeMessage(mode)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *SetGameCodeProcessor) setGameCodeMessage(mode int) error {
	wsCollection := sockets.NewCollection("priceIsRight")
	wsCollection.SendText("CODE "+strconv.Itoa(mode), wsCollection.All()...)

	return nil
}
