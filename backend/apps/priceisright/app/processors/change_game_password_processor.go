package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangeGamePasswordProcessorInterface interface {
	ChangeGamePassword(value string) (*models.Game, error)
}

type ChangeGamePasswordProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeGamePasswordProcessor) New(gameRepository repositories.GameRepository) *ChangeGamePasswordProcessor {
	return &ChangeGamePasswordProcessor{
		gameRepository,
	}
}

func (P *ChangeGamePasswordProcessor) ChangeGamePassword(value string) (*models.Game, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	game.Password = value
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.GamePasswordChangedMessage(value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangeGamePasswordProcessor) GamePasswordChangedMessage(value string) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Password string
	}{
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"GAME_PASSWORD_CHANGED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
