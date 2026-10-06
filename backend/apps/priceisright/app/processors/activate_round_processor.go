package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ActivateRoundProcessorInterface interface {
	ActivateRound(round int) (*models.Round, error)
}

type ActivateRoundProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ActivateRoundProcessor) New(gameRepository repositories.GameRepository) *ActivateRoundProcessor {
	return &ActivateRoundProcessor{
		gameRepository,
	}
}

func (P *ActivateRoundProcessor) ActivateRound(round int) (*models.Round, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	if round < 0 || round >= len(game.Rounds) {
		round = game.RoundNumber
	} else {
		game.RoundNumber = round
	}
	game.CurrentRound = game.Rounds[game.RoundNumber]
	game.CurrentRound.Activated = time.Now()
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.ActivateRoundMessage(*game.CurrentRound)
	if err != nil {
		return game.CurrentRound, err
	}

	return game.CurrentRound, nil
}

func (P *ActivateRoundProcessor) ActivateRoundMessage(round models.Round) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Round models.Round
	}{
		round,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"ROUND_ACTIVATED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_AUDIENCE))

	data.Round.Value = 0
	data.Round.Guesses = make([]*models.Guess, 0)
	msg, err = json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"ROUND_ACTIVATED",
		data,
	})

	wsCollection.SendText(string(msg), socks...)

	return nil
}
