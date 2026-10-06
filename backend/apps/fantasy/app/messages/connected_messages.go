package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

// ConnectedMessages brings a newly connected socket up to date, without disturbing any
// other sockets the same player has open. Must be called with the game locked.
func ConnectedMessages(game *models.Game, player *models.Player, connection *sockets.Connection) error {
	msg, err := marshal("GAME_UPDATE", View(game, player))
	if err != nil {
		return err
	}

	sockets.GetCollection(game.Key).SendText(msg, connection)
	return nil
}
