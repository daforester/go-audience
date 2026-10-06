package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type NewGuessProcessorInterface interface {
	NewGuess(player *models.Player, value int64) (*models.Guess, error)
}

type NewGuessProcessor struct {
	gameRepository repositories.GameRepository
}

func (P NewGuessProcessor) New(gameRepository repositories.GameRepository) *NewGuessProcessor {
	return &NewGuessProcessor{
		gameRepository,
	}
}

func (P *NewGuessProcessor) NewGuess(player *models.Player, value int64) (*models.Guess, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	if game.RoundNumber < 0 || game.RoundNumber > len(game.Rounds)-1 {
		return nil, errors.New("no round active")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	guess := &models.Guess{
		PlayerUUID: player.UUID,
		Player:     player,
		RoundUUID:  game.CurrentRound.UUID,
		Round:      game.CurrentRound,
		Value:      value,
		Timestamp:  time.Now(),
	}

	game.CurrentRound.Guesses = append(game.CurrentRound.Guesses, guess)

	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.newGuessMessage(guess)
	if err != nil {
		return guess, err
	}

	return guess, nil
}

func (P *NewGuessProcessor) newGuessMessage(guess *models.Guess) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Guess *models.Guess
	}{
		guess,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"GUESS_ADDED",
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
