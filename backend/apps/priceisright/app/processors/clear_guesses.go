package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ClearGuessesProcessorInterface interface {
	ClearGuesses(round int) ([]*models.Guess, error)
}

type ClearGuessesProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ClearGuessesProcessor) New(gameRepository repositories.GameRepository) *ClearGuessesProcessor {
	return &ClearGuessesProcessor{
		gameRepository,
	}
}

func (P *ClearGuessesProcessor) ClearGuesses(round int) ([]*models.Guess, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(round) >= len(game.Rounds) {
		return nil, errors.New("index out of range")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	game.Rounds[round].Guesses = make([]*models.Guess, 0)
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.clearGuessesMessage(round)
	if err != nil {
		return game.Rounds[round].Guesses, err
	}

	return game.Rounds[round].Guesses, nil
}

func (P *ClearGuessesProcessor) clearGuessesMessage(round int) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Round int
	}{
		round,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CLEAR_GUESSES",
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
