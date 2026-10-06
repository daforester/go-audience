package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangeGuessProcessorInterface interface {
	ChangeGuess(round uint, index uint, value int64) (*models.Game, error)
}

type ChangeGuessProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeGuessProcessor) New(gameRepository repositories.GameRepository) *ChangeGuessProcessor {
	return &ChangeGuessProcessor{
		gameRepository,
	}
}

func (P *ChangeGuessProcessor) ChangeGuess(round uint, index uint, value int64) (*models.Game, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(index) >= len(game.Rounds[round].Guesses) {
		return nil, errors.New("index out of range")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	game.Rounds[round].Guesses[index].Value = value
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.guessChangedMessage(round, index, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangeGuessProcessor) guessChangedMessage(round uint, index uint, value int64) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Round uint
		Index uint
		Value int64
	}{
		round,
		index,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"GUESS_CHANGED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
