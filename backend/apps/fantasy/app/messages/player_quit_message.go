package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

func PlayerQuitMessage(player *models.Player) error {
	return sendToClasses(player.RoomKey, "PLAYER_QUIT", player.Copy(), models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_DISPLAY)
}
