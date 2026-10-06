package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type AssignNameProcessorInterface interface {
	AssignName(roomKey string, playerUUID string, name string) (*models.Player, error)
}

type AssignNameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AssignNameProcessor) New(gameRepository repositories.GameRepository) *AssignNameProcessor {
	return &AssignNameProcessor{
		gameRepository,
	}
}

func (P *AssignNameProcessor) AssignName(roomKey string, playerUUID string, name string) (*models.Player, error) {
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

	player.CharacterName = name

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return player, err
	}

	err = P.playerUpdateMessage(player)
	if err != nil {
		return player, err
	}

	return player, nil
}

func (P *AssignNameProcessor) playerUpdateMessage(player *models.Player) error {
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
