package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangePlayerNameProcessorInterface interface {
	ChangePlayerName(roomKey string, index uint, value string) (*models.Game, error)
}

type ChangePlayerNameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangePlayerNameProcessor) New(gameRepository repositories.GameRepository) *ChangePlayerNameProcessor {
	return &ChangePlayerNameProcessor{
		gameRepository,
	}
}

func (P *ChangePlayerNameProcessor) ChangePlayerName(roomKey string, index uint, value string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(index) >= len(game.Players) {
		return nil, errors.New("index out of range")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	game.Players[index].Name = value
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.playerNameChangedMessage(roomKey, game.Players[index], index, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangePlayerNameProcessor) playerNameChangedMessage(roomKey string, player *models.Player, index uint, value string) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		Index uint
		Name  string
	}{
		index,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"PLAYER_NAME_CHANGED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasUUID(player.UUID)
	wsCollection.SendText(string(msg), socks...)

	return nil
}
