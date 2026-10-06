package messages

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"encoding/json"
	"errors"
)

func HealthUpdateMessage(roomKey string, player *models.Player, index uint, health int) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"HEALTH_UPDATED",
		map[string]interface{}{
			"Health": health,
		},
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasUUID(player.UUID)
	wsCollection.SendText(string(msg), socks...)

	msg, err = json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"PLAYER_HEALTH_CHANGED",
		map[string]interface{}{
			"Index":  index,
			"Health": health,
		},
	})

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
