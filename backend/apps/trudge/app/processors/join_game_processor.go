package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type JoinGameProcessorInterface interface {
	JoinGame(playerName string, roomCode string) (*models.Player, error)
}

type JoinGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P JoinGameProcessor) New(gameRepository repositories.GameRepository) *JoinGameProcessor {
	return &JoinGameProcessor{
		gameRepository,
	}
}

func (P *JoinGameProcessor) JoinGame(playerName string, roomCode string) (*models.Player, error) {
	player := new(models.Player)
	player.Name = playerName
	player.RoomKey = strings.ToUpper(roomCode)
	player.Class = models.PLAYER_CLASS_AUDIENCE
	player.Health = 5
	player.Experience = 0
	player.Vote = -1
	player.Option = -1
	player.Equipment = make([]string, 0)

	id := uuid.New()
	player.UUID = id.String()

	game := P.gameRepository.GetByKey(player.RoomKey)
	if game == nil {
		return nil, errors.New("no game with room key")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Updated = time.Now()
	game.AddPlayer(player)
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, err
	}

	err = P.newPlayerMessage(player)

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (P *JoinGameProcessor) newPlayerMessage(player *models.Player) error {

	msg, err := json.Marshal(struct {
		Command string
		Data    *models.Player
	}{
		"NEW_PLAYER",
		player,
	})
	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	wsCollection := sockets.NewCollection(player.RoomKey)

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)
	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)
	return nil
}
