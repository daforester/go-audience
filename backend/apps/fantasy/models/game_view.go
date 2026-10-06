package models

// GameView is what a player is allowed to see of a game. It is always built from copies,
// so it can be marshalled after the game lock has been released.
type GameView struct {
	Key      string
	GameCode int
	Players  []*Player `json:",omitempty"`
	Volume   int
}

func (g *Game) ViewFor(player *Player) *GameView {
	view := &GameView{
		Key:      g.Key,
		GameCode: g.GameCode,
		Volume:   g.Volume,
	}

	if player.Class == PLAYER_CLASS_ADMIN || player.Class == PLAYER_CLASS_DISPLAY {
		view.Players = g.CopyPlayers()
	}

	return view
}
