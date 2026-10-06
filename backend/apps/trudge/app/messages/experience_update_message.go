package messages

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"encoding/json"
	"errors"
)

func ExperienceUpdateMessage(roomKey string, player *models.Player, index uint, experience int) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"EXPERIENCE_UPDATED",
		map[string]interface{}{
			"Experience": experience,
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
		"PLAYER_EXPERIENCE_CHANGED",
		map[string]interface{}{
			"Index":      index,
			"Experience": experience,
		},
	})

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
