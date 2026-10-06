package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type AddCharacterEquipmentProcessorInterface interface {
	AddCharacterEquipment(roomKey string, playerIndex uint, value string) (*models.Game, error)
}

type AddCharacterEquipmentProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AddCharacterEquipmentProcessor) New(gameRepository repositories.GameRepository) *AddCharacterEquipmentProcessor {
	return &AddCharacterEquipmentProcessor{
		gameRepository,
	}
}

func (P *AddCharacterEquipmentProcessor) AddCharacterEquipment(roomKey string, playerIndex uint, value string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(playerIndex) >= len(game.Players) {
		return nil, errors.New("player index out of range")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.Players[playerIndex]

	if player.Equipment == nil {
		player.Equipment = make([]string, 0)
	}
	player.Equipment = append(player.Equipment, value)

	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.characterEquipmentAddedMessage(roomKey, game.Players[playerIndex], playerIndex, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *AddCharacterEquipmentProcessor) characterEquipmentAddedMessage(roomKey string, player *models.Player, playerIndex uint, value string) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		UUID        string
		PlayerIndex uint
		Equipment   string
	}{
		player.UUID,
		playerIndex,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CHARACTER_EQUIPMENT_ADDED",
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
