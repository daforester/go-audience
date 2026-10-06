package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

func PlayerUpdateMessage(player *models.Player) error {
	err := sendToClasses(player.RoomKey, "PLAYER_UPDATE", player.Copy(), models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_DISPLAY)
	if err != nil {
		return err
	}

	if player.Class == models.PLAYER_CLASS_ADMIN || player.Class == models.PLAYER_CLASS_DISPLAY {
		return nil
	}

	return sendToPlayer(player.RoomKey, player.UUID, "PLAYER_UPDATE", player.Copy())
}
