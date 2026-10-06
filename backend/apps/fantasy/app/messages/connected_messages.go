package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

// ConnectedMessages brings a newly connected player's sockets up to date. Must be called
// with the game locked.
func ConnectedMessages(game *models.Game, player *models.Player) error {
	if player.Class == models.PLAYER_CLASS_ADMIN || player.Class == models.PLAYER_CLASS_DISPLAY {
		err := sendToPlayer(game.Key, player.UUID, "SET_VOLUME", map[string]interface{}{
			"Volume": game.Volume,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
