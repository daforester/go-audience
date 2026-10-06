package models

import "time"

type Game struct {
	Key                string
	GameCode           int
	Quest              string
	Players            []*Player
	SelectedOption     int
	Prompt             string
	Options            []string
	Votes              []map[string]uint
	SuggestionsEnabled bool
	Volume             int
	Created            time.Time
	Updated            time.Time
	Deleted            *time.Time
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

func (g *Game) GetAudience() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_AUDIENCE)
}

func (g *Game) GetGods() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_GOD)
}

func (g *Game) GetLeader() *Player {
	return g.GetPlayerByType(PLAYER_CLASS_LEADER)
}

func (g *Game) GetBard() *Player {
	return g.GetPlayerByType(PLAYER_CLASS_BARD)
}

func (g *Game) GetWarrior() *Player {
	return g.GetPlayerByType(PLAYER_CLASS_WARRIOR)
}

func (g *Game) GetWizard() *Player {
	return g.GetPlayerByType(PLAYER_CLASS_WIZARD)
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
