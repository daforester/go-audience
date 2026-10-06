package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ClearRequestOptionsProcessorInterface interface {
	ClearRequestOptions(roomKey string, prompt string) (*models.Game, error)
}

type ClearRequestOptionsProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ClearRequestOptionsProcessor) New(gameRepository repositories.GameRepository) *ClearRequestOptionsProcessor {
	return &ClearRequestOptionsProcessor{
		gameRepository,
	}
}

func (P *ClearRequestOptionsProcessor) ClearRequestOptions(roomKey string, prompt string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Prompt = prompt
	game.Options = make([]string, 0)
	game.Votes = make([]map[string]uint, 0)
	game.Updated = time.Now()

	for _, player := range game.Players {
		player.Vote = -1
		player.Option = -1
	}

	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.optionsClearedMessage(roomKey)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ClearRequestOptionsProcessor) optionsClearedMessage(roomKey string) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
	}{
		"OPTIONS_CLEARED",
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	wsCollection.SendText(string(msg))

	return nil
}
