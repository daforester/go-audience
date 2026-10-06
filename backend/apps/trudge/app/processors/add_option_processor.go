package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type AddOptionProcessorInterface interface {
	AddOption(roomKey string, playerUUID string, option string) (*models.Game, error)
}

type AddOptionProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AddOptionProcessor) New(gameRepository repositories.GameRepository) *AddOptionProcessor {
	return &AddOptionProcessor{
		gameRepository,
	}
}

func (P *AddOptionProcessor) AddOption(roomKey string, playerUUID string, option string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	for _, player := range game.Players {
		if player.UUID == playerUUID {
			player.Option = len(game.Options)
			break
		}
	}
	game.Options = append(game.Options, option)
	game.Votes = append(game.Votes, make(map[string]uint))

	i := len(game.Votes) - 1
	game.Votes[i]["Total"] = 0
	game.Votes[i][string(models.PLAYER_CLASS_AUDIENCE)] = 0
	game.Votes[i][string(models.PLAYER_CLASS_BARD)] = 0
	game.Votes[i][string(models.PLAYER_CLASS_GOD)] = 0
	game.Votes[i][string(models.PLAYER_CLASS_WARRIOR)] = 0
	game.Votes[i][string(models.PLAYER_CLASS_WIZARD)] = 0

	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.optionAddedMessage(game.Key, option)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *AddOptionProcessor) optionAddedMessage(roomKey string, option string) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Option  string
	}{
		"OPTION_ADDED",
		option,
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
