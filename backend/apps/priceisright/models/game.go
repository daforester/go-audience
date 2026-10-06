package models

import (
	"time"
)

type Game struct {
	GameCode     int
	Password     string
	Players      []*Player
	Volume       int
	RoundNumber  int
	CurrentRound *Round
	Rounds       []*Round
	BadgeCheck   bool

	Created time.Time
	Updated time.Time
	Deleted *time.Time
}

func (g Game) New() *Game {
	game := new(Game)
	game.Players = make([]*Player, 0)
	game.Rounds = make([]*Round, 0)
	game.BadgeCheck = false
	return game
}

func (g *Game) ComingSoon() {
	g.GameCode = 1
}

func (g *Game) StartingSoon() {
	g.GameCode = 2
}

func (g *Game) PleaseWait() {
	g.GameCode = 10
}

func (g *Game) AllowGuesses() {
	g.GameCode = 20
}

func (g *Game) StopGuesses() {
	g.GameCode = 30
}

func (g *Game) ShowResult() {
	g.GameCode = 40
}

func (g *Game) GameOver() {
	g.GameCode = 100
}

func (g *Game) NextRound() {
	g.RoundNumber = g.RoundNumber + 1
	if g.RoundNumber > len(g.Rounds) {
		g.RoundNumber = 0
		g.CurrentRound = nil
		return
	}
	for _, v := range g.Rounds {
		if v.Order == g.RoundNumber {
			g.CurrentRound = v
			return
		}
	}
}

func (g *Game) AddPlayer(player *Player) []*Player {
	g.Players = append(g.Players, player)
	return g.Players
}

func (g *Game) RemovePlayer(player *Player) {
	for i, p := range g.Players {
		if p.UUID == player.UUID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
		}
	}
}

func (g *Game) ToggleBadgeCheck() {
	g.BadgeCheck = !g.BadgeCheck
}

func (g *Game) GetAudience() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_AUDIENCE)
}

func (g *Game) GetLeader() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_LEADER)
}

func (g *Game) GetPlayerByUUID(uuid string) *Player {
	for _, player := range g.Players {
		if player.UUID == uuid {
			return player
		}
	}
	return nil
}

func (g *Game) GetPlayerByType(playerClass PlayerClass) *Player {
	for _, player := range g.Players {
		if player.Class == playerClass {
			return player
		}
	}
	return nil
}

func (g *Game) GetPlayersByType(playerClass PlayerClass) []*Player {
	var players []*Player
	for _, player := range g.Players {
		if player.Class == playerClass {
			players = append(players, player)
		}
	}
	return players
}
