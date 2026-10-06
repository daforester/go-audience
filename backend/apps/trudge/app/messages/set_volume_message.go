package messages

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"encoding/json"
	"errors"
)

func SetVolumeMessage(roomKey string, volume int) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"SET_VOLUME",
		map[string]interface{}{
			"Volume": volume,
		},
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
