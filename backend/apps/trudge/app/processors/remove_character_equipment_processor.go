package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type RemoveCharacterEquipmentProcessorInterface interface {
	RemoveCharacterEquipment(roomKey string, playerIndex uint, equipmentIndex uint) (*models.Game, error)
}

type RemoveCharacterEquipmentProcessor struct {
	gameRepository repositories.GameRepository
}

func (P RemoveCharacterEquipmentProcessor) New(gameRepository repositories.GameRepository) *RemoveCharacterEquipmentProcessor {
	return &RemoveCharacterEquipmentProcessor{
		gameRepository,
	}
}

func (P *RemoveCharacterEquipmentProcessor) RemoveCharacterEquipment(roomKey string, playerIndex uint, equipmentIndex uint) (*models.Game, error) {
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

	player := game.Players[playerIndex]
	player.Equipment = append(player.Equipment[:equipmentIndex], player.Equipment[equipmentIndex+1:]...)

	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.characterEquipmentRemovedMessage(roomKey, game.Players[playerIndex], playerIndex, equipmentIndex)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *RemoveCharacterEquipmentProcessor) characterEquipmentRemovedMessage(roomKey string, player *models.Player, playerIndex uint, equipmentIndex uint) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		UUID           string
		PlayerIndex    uint
		EquipmentIndex uint
	}{
		player.UUID,
		playerIndex,
		equipmentIndex,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"CHARACTER_EQUIPMENT_REMOVED",
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
