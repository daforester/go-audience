package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangeCharacterRaceProcessorInterface interface {
	ChangeCharacterRace(roomKey string, index uint, value string) (*models.Game, error)
}

type ChangeCharacterRaceProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeCharacterRaceProcessor) New(gameRepository repositories.GameRepository) *ChangeCharacterRaceProcessor {
	return &ChangeCharacterRaceProcessor{
		gameRepository,
	}
}

func (P *ChangeCharacterRaceProcessor) ChangeCharacterRace(roomKey string, index uint, value string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(index) >= len(game.Players) {
		return nil, errors.New("index out of range")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	game.Players[index].CharacterRace = value
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.characterRaceChangedMessage(roomKey, game.Players[index], index, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangeCharacterRaceProcessor) characterRaceChangedMessage(roomKey string, player *models.Player, index uint, value string) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		UUID          string
		Index         uint
		CharacterRace string
	}{
		player.UUID,
		index,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CHARACTER_RACE_CHANGED",
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
