package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

type JoinGameProcessorInterface interface {
	JoinGame(playerName string, badgeNumber string, password string) (*models.Player, error)
}

type JoinGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P JoinGameProcessor) New(gameRepository repositories.GameRepository) *JoinGameProcessor {
	return &JoinGameProcessor{
		gameRepository,
	}
}

func (P *JoinGameProcessor) JoinGame(playerName string, badgeNumber string, password string) (*models.Player, error) {
	player := new(models.Player)
	player.Name = playerName
	player.BadgeNumber = badgeNumber
	player.Class = models.PLAYER_CLASS_AUDIENCE

	id := uuid.New()
	player.UUID = id.String()

	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("no game found")
	}

	if game.BadgeCheck && helpers.ValidBadgeNumber(badgeNumber) == false {
		return nil, errors.New("INVALID_BADGE")
	}

	if game.Password == password {
		player.Class = models.PLAYER_CLASS_LEADER
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()
	game.Updated = time.Now()
	game.AddPlayer(player)
	P.gameRepository.Save(game)

	err := P.newPlayerMessage(player)

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

	wsCollection := sockets.NewCollection("priceIsRight")

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)
	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)
	return nil
}
