package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type SaveRoundsProcessorInterface interface {
	SaveRounds(rounds []*models.Round) ([]*models.Round, error)
}

type SaveRoundsProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SaveRoundsProcessor) New(gameRepository repositories.GameRepository) *SaveRoundsProcessor {
	return &SaveRoundsProcessor{
		gameRepository,
	}
}

func (P *SaveRoundsProcessor) SaveRounds(rounds []*models.Round) ([]*models.Round, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	// Copy existing data
	for i, newRound := range rounds {
		newRound.Order = i + 1
		for _, oldRound := range game.Rounds {
			if oldRound == nil {
				continue
			}
			if oldRound.UUID == newRound.UUID {
				if newRound.Image == "" {
					newRound.Image = oldRound.Image
				}
				newRound.Guesses = oldRound.Guesses
			}
			if newRound.Guesses == nil {
				newRound.Guesses = make([]*models.Guess, 0)
			}
		}
	}

	game.Rounds = rounds

	if game.RoundNumber < len(game.Rounds) {
		game.CurrentRound = game.Rounds[game.RoundNumber]
	} else {
		game.RoundNumber = 0
		if game.RoundNumber < len(game.Rounds) {
			game.CurrentRound = game.Rounds[game.RoundNumber]
		} else {
			game.CurrentRound = nil
		}
	}

	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.roundsUpdatedMessage(rounds)
	if err != nil {
		return game.Rounds, err
	}

	return game.Rounds, nil
}

func (P *SaveRoundsProcessor) roundsUpdatedMessage(rounds []*models.Round) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Rounds []*models.Round
	}{
		rounds,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"ROUNDS_UPDATED",
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
