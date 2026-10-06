package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

func NewPlayerMessage(player *models.Player) error {
	return sendToClasses(player.RoomKey, "NEW_PLAYER", player.Copy(), models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_DISPLAY)
}
