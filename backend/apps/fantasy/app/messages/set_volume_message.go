package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

func SetVolumeMessage(roomKey string, volume int) error {
	return sendToClasses(roomKey, "SET_VOLUME", map[string]interface{}{
		"Volume": volume,
	}, models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_DISPLAY)
}
