package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangePlayerBadgeProcessorInterface interface {
	ChangePlayerBadge(index uint, value string) (*models.Game, error)
}

type ChangePlayerBadgeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangePlayerBadgeProcessor) New(gameRepository repositories.GameRepository) *ChangePlayerBadgeProcessor {
	return &ChangePlayerBadgeProcessor{
		gameRepository,
	}
}

func (P *ChangePlayerBadgeProcessor) ChangePlayerBadge(index uint, value string) (*models.Game, error) {
	game := P.gameRepository.Get()
	if game == nil {
		return nil, errors.New("game not found")
	}
	if int(index) >= len(game.Players) {
		return nil, errors.New("index out of range")
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()

	game.Players[index].BadgeNumber = value
	game.Updated = time.Now()
	P.gameRepository.Save(game)

	err := P.playerBadgeChangedMessage(game.Players[index], index, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangePlayerBadgeProcessor) playerBadgeChangedMessage(player *models.Player, index uint, value string) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Index       uint
		BadgeNumber string
	}{
		index,
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"PLAYER_BADGE_CHANGED",
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
