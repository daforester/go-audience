package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type AssignRoleProcessorInterface interface {
	AssignPlayer(roomKey string, playerUUID string, class models.PlayerClass) (*models.Player, error)
}

type AssignRoleProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AssignRoleProcessor) New(gameRepository repositories.GameRepository) *AssignRoleProcessor {
	return &AssignRoleProcessor{
		gameRepository,
	}
}

func (P *AssignRoleProcessor) AssignPlayer(roomKey string, playerUUID string, class models.PlayerClass) (*models.Player, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	var player *models.Player
	for _, p := range game.Players {
		if p.UUID == playerUUID {
			player = p
		}
	}

	if player == nil {
		return nil, errors.New("player not found")
	}

	player.Class = class

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return player, err
	}

	err = P.assignedPlayerMessage(player)
	if err != nil {
		return player, err
	}

	err = P.swapSocketTags(player)
	if err != nil {
		return player, err
	}

	return player, nil
}

func (P *AssignRoleProcessor) swapSocketTags(player *models.Player) error {
	wsCollection := sockets.NewCollection(player.RoomKey)
	socks := wsCollection.HasUUID(player.UUID)
	for _, sock := range socks {
		_, err := sock.SwapTags(models.TrudgePlayerClassesStrings, []string{string(player.Class)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (P *AssignRoleProcessor) assignedPlayerMessage(player *models.Player) error {
	wsCollection := sockets.NewCollection(player.RoomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    *models.Player
	}{
		"PLAYER_UPDATE",
		player,
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
