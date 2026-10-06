package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
)

type DiceRollProcessorInterface interface {
	RollDice(roomKey string) (*models.Game, error)
}

type DiceRollProcessor struct {
	gameRepository repositories.GameRepository
}

func (P DiceRollProcessor) New(gameRepository repositories.GameRepository) *DiceRollProcessor {
	return &DiceRollProcessor{
		gameRepository,
	}
}

func (P *DiceRollProcessor) RollDice(roomKey string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	err := P.rollDiceMessage(game.Key)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *DiceRollProcessor) rollDiceMessage(roomKey string) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
	}{
		"ROLL_DICE",
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_BARD))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_WARRIOR))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_WIZARD))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
