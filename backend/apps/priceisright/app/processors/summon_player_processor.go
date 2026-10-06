package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
)

type SummonPlayerProcessorInterface interface {
	SummonPlayer(player *models.Player, value int) error
}

type SummonPlayerProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SummonPlayerProcessor) New(gameRepository repositories.GameRepository) *SummonPlayerProcessor {
	return &SummonPlayerProcessor{
		gameRepository,
	}
}

func (P *SummonPlayerProcessor) SummonPlayer(player *models.Player, value int) error {
	return P.SummonPlayerMessage(player, value)
}

func (P *SummonPlayerProcessor) SummonPlayerMessage(player *models.Player, value int) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		UUID        string
		Name        string
		BadgeNumber string
		Value       int
	}{
		"",
		"",
		"",
		value,
	}

	if player != nil {
		data.UUID = player.UUID
		data.Name = player.Name
		data.BadgeNumber = player.BadgeNumber
		data.Value = value
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"COME_ON_DOWN",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	if player != nil {
		socks = wsCollection.HasUUID(player.UUID)
		wsCollection.SendText(string(msg), socks...)
	} else {
		socks = wsCollection.HasTag(string(models.PLAYER_CLASS_AUDIENCE))
		wsCollection.SendText(string(msg), socks...)
	}

	return nil
}
