package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type SetRoundProcessorInterface interface {
	SetRound(round int) (*models.Game, error)
}

type SetRoundProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SetRoundProcessor) New(gameRepository repositories.GameRepository) *SetRoundProcessor {
	return &SetRoundProcessor{
		gameRepository,
	}
}

func (P *SetRoundProcessor) SetRound(round int) (*models.Game, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	if round < 0 || round >= len(game.Rounds) {
		return nil, errors.New("round out of range")
	}
	game.RoundNumber = round
	game.CurrentRound = game.Rounds[game.RoundNumber]
	game.CurrentRound.Activated = time.Now()
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.setRoundMessage(game.RoundNumber, game.CurrentRound)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *SetRoundProcessor) setRoundMessage(roundNumber int, round *models.Round) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		RoundNumber int
		Round       *models.Round
	}{
		roundNumber,
		round,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CURRENT_ROUND",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	audienceRound := new(models.Round)
	audienceRound.Name = data.Round.Name
	audienceRound.UUID = data.Round.UUID
	audienceRound.Image = data.Round.Image
	audienceRound.Order = data.Round.Order

	data.Round = audienceRound

	msg, err = json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CURRENT_ROUND",
		data,
	})

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_AUDIENCE))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
