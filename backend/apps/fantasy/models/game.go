package models

import "time"

type Game struct {
	Key      string
	GameCode int
	Players  []*Player
	Teams    []*Team
	Volume   int
	Created  time.Time
	Updated  time.Time
	Deleted  *time.Time
}

func NewGame() *Game {
	game := new(Game)
	game.GameCode = GAME_CODE_LOBBY
	game.Players = make([]*Player, 0)
	game.Teams = make([]*Team, 0)
	game.Volume = 50
	return game
}

func (g *Game) AddPlayer(player *Player) []*Player {
	g.Players = append(g.Players, player)
	return g.Players
}

func (g *Game) RemovePlayer(player *Player) {
	for i, p := range g.Players {
		if p.UUID == player.UUID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			return
		}
	}
}

func (g *Game) GetAdmins() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_ADMIN)
}

func (g *Game) GetAudience() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_AUDIENCE)
}

func (g *Game) GetDirectors() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_DIRECTOR)
}

func (g *Game) GetStaff() []*Player {
	return g.GetPlayersByType(PLAYER_CLASS_STAFF)
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

// CopyPlayers returns copies of every player, safe to use after the game lock is released
func (g *Game) CopyPlayers() []*Player {
	players := make([]*Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p.Copy())
	}
	return players
}

func (g *Game) GetTeam(id string) *Team {
	if id == "" {
		return nil
	}
	for _, team := range g.Teams {
		if team.ID == id {
			return team
		}
	}
	return nil
}

func (g *Game) GetTeamDirectors(teamID string) []*Player {
	var players []*Player
	for _, player := range g.GetDirectors() {
		if player.TeamID == teamID {
			players = append(players, player)
		}
	}
	return players
}

func (g *Game) CopyTeams() []*Team {
	teams := make([]*Team, 0, len(g.Teams))
	for _, t := range g.Teams {
		teams = append(teams, t.Copy())
	}
	return teams
}
