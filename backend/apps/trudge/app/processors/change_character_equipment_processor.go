package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangeCharacterEquipmentProcessorInterface interface {
	ChangeCharacterEquipment(roomKey string, playerIndex uint, equipmentIndex uint, value string) (*models.Game, error)
}

type ChangeCharacterEquipmentProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeCharacterEquipmentProcessor) New(gameRepository repositories.GameRepository) *ChangeCharacterEquipmentProcessor {
	return &ChangeCharacterEquipmentProcessor{
		gameRepository,
	}
}

func (P *ChangeCharacterEquipmentProcessor) ChangeCharacterEquipment(roomKey string, playerIndex uint, equipmentIndex uint, value string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(playerIndex) >= len(game.Players) {
		return nil, errors.New("player index out of range")
	}
	if int(equipmentIndex) >= len(game.Players[playerIndex].Equipment) {
		return nil, errors.New("equipment index out of range")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	game.Players[playerIndex].Equipment[equipmentIndex] = value
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.characterEquipmentChangedMessage(roomKey, game.Players[playerIndex], playerIndex, equipmentIndex, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangeCharacterEquipmentProcessor) characterEquipmentChangedMessage(roomKey string, player *models.Player, playerIndex uint, equipmentIndex uint, value string) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		UUID           string
		PlayerIndex    uint
		EquipmentIndex uint
		Equipment      string
	}{
		player.UUID,
		playerIndex,
		equipmentIndex,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CHARACTER_EQUIPMENT_CHANGED",
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
