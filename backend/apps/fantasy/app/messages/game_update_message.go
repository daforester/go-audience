package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"github.com/spf13/viper"
)

// View builds the view a player may see, with server settings the admin needs
func View(game *models.Game, player *models.Player) *models.GameView {
	view := game.ViewFor(player)
	if player.Class == models.PLAYER_CLASS_ADMIN {
		view.MaxTeams = viper.GetInt("MAX_TEAMS")
	}
	return view
}

// GameUpdateMessage sends every player their own view of the game. Used when the game
// moves on or changes shape; frequent small changes use their own messages instead.
// Must be called with the game locked.
func GameUpdateMessage(game *models.Game) error {
	for _, player := range game.Players {
		err := sendToPlayer(game.Key, player.UUID, "GAME_UPDATE", View(game, player))
		if err != nil {
			return err
		}
	}
	return nil
}

// PlayerGameUpdateMessage sends one player their view of the game
func PlayerGameUpdateMessage(game *models.Game, player *models.Player) error {
	return sendToPlayer(game.Key, player.UUID, "GAME_UPDATE", View(game, player))
}
