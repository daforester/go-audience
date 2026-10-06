package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
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
	playerName = CleanName(playerName)
	if playerName == "" {
		return nil, fmt.Errorf("%w: name required", ErrInvalid)
	}

	player := new(models.Player)
	player.Name = playerName
	player.RoomKey = strings.ToUpper(strings.TrimSpace(roomCode))
	player.Class = models.PLAYER_CLASS_AUDIENCE

	id := uuid.New()
	player.UUID = id.String()

	game := P.gameRepository.GetByKey(player.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Updated = time.Now()
	game.AddPlayer(player)
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, err
	}

	err = messages.NewPlayerMessage(player)

	if err != nil {
		return nil, err
	}

	return player.Copy(), nil
}
