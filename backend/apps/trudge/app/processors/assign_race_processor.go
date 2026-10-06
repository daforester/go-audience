package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type AssignRaceProcessorInterface interface {
	AssignRace(roomKey string, playerUUID string, race string) (*models.Player, error)
}

type AssignRaceProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AssignRaceProcessor) New(gameRepository repositories.GameRepository) *AssignRaceProcessor {
	return &AssignRaceProcessor{
		gameRepository,
	}
}

func (P *AssignRaceProcessor) AssignRace(roomKey string, playerUUID string, race string) (*models.Player, error) {
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

	player.CharacterRace = race

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

func (P *AssignRaceProcessor) playerUpdateMessage(player *models.Player) error {
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
