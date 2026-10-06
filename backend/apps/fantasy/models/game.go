package models

import (
	"errors"
	"sort"
	"time"
)

type Game struct {
	Key                string
	GameCode           int
	Players            []*Player
	Teams              []*Team
	Submissions        []*Submission
	Rounds             []*Round
	SubmissionsEnabled bool
	ShowSubmitterNames bool
	PinnedTopicID      string
	Volume             int
	LastSubmission     map[string]time.Time `json:"-"`
	Created            time.Time
	Updated            time.Time
	Deleted            *time.Time
}

func NewGame() *Game {
	game := new(Game)
	game.GameCode = GAME_CODE_LOBBY
	game.Players = make([]*Player, 0)
	game.Teams = make([]*Team, 0)
	game.Submissions = make([]*Submission, 0)
	game.Rounds = make([]*Round, 0)
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

func (g *Game) AddTeam(team *Team) {
	g.Teams = append(g.Teams, team)
}

// RemoveTeam removes the team, its directors return to the audience
func (g *Game) RemoveTeam(id string) {
	for i, t := range g.Teams {
		if t.ID == id {
			g.Teams = append(g.Teams[:i], g.Teams[i+1:]...)
			break
		}
	}
	for _, p := range g.GetTeamDirectors(id) {
		p.Class = PLAYER_CLASS_AUDIENCE
		p.TeamID = ""
	}
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

func (g *Game) GetSubmission(id string) *Submission {
	for _, s := range g.Submissions {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (g *Game) AddSubmission(submission *Submission) {
	g.Submissions = append(g.Submissions, submission)
}

// TopPicks are the approved submissions in the order they would be played: most staff
// votes first, oldest first on a tie. At most n are returned, n <= 0 returns them all.
func (g *Game) TopPicks(n int) []*Submission {
	picks := make([]*Submission, 0)
	for _, s := range g.Submissions {
		if s.Status == SUBMISSION_STATUS_APPROVED {
			picks = append(picks, s)
		}
	}
	sort.SliceStable(picks, func(i, j int) bool {
		a, b := picks[i], picks[j]
		if a.StaffVoteCount() != b.StaffVoteCount() {
			return a.StaffVoteCount() > b.StaffVoteCount()
		}
		return a.Created.Before(b.Created)
	})
	if n > 0 && len(picks) > n {
		picks = picks[:n]
	}
	return picks
}

// PinnedTopic is the submission the admin has chosen to play next, if it can still be played
func (g *Game) PinnedTopic() *Submission {
	if g.PinnedTopicID == "" {
		return nil
	}
	s := g.GetSubmission(g.PinnedTopicID)
	if s == nil || !s.Eligible() {
		return nil
	}
	return s
}

// NextTopic is the admin's pinned submission if there is one, otherwise the staff's
// top pick
func (g *Game) NextTopic() *Submission {
	if pinned := g.PinnedTopic(); pinned != nil {
		return pinned
	}
	if picks := g.TopPicks(1); len(picks) > 0 {
		return picks[0]
	}
	return nil
}

func (g *Game) GetCurrentRound() *Round {
	if len(g.Rounds) == 0 {
		return nil
	}
	return g.Rounds[len(g.Rounds)-1]
}

// RoundInProgress reports whether a round has started and its result is not yet shown
func (g *Game) RoundInProgress() bool {
	return g.GameCode == GAME_CODE_TOPIC_REVEAL ||
		g.GameCode == GAME_CODE_RESPONSES ||
		g.GameCode == GAME_CODE_AUDIENCE_VOTE
}

// ReadyForRound checks there are at least two teams and every team has a director
func (g *Game) ReadyForRound() error {
	if len(g.Teams) < 2 {
		return errors.New("at least two director teams are needed")
	}
	for _, t := range g.Teams {
		if len(g.GetTeamDirectors(t.ID)) == 0 {
			return errors.New("team " + t.Name + " has no directors")
		}
	}
	return nil
}

// StartRound plays the submission as the next topic. Teams present in an order that
// rotates each round, so no team always goes first.
func (g *Game) StartRound(submission *Submission) *Round {
	number := len(g.Rounds) + 1

	order := make([]string, 0, len(g.Teams))
	for i := range g.Teams {
		order = append(order, g.Teams[(i+number-1)%len(g.Teams)].ID)
	}

	now := time.Now()
	submission.Status = SUBMISSION_STATUS_PLAYED
	submission.Played = &now
	if g.PinnedTopicID == submission.ID {
		g.PinnedTopicID = ""
	}

	round := NewRound(number, submission.ID, order)
	g.Rounds = append(g.Rounds, round)
	return round
}

// DecideRound awards the round to the teams with the most votes. Only counts once.
func (g *Game) DecideRound() *Round {
	round := g.GetCurrentRound()
	if round == nil || round.Decided {
		return round
	}

	round.Decided = true
	round.WinnerTeamIDs = round.Winners()
	round.PresentingTeam = ""

	for teamID, votes := range round.Tally() {
		if team := g.GetTeam(teamID); team != nil {
			team.TotalVotes += votes
		}
	}
	for _, teamID := range round.WinnerTeamIDs {
		if team := g.GetTeam(teamID); team != nil {
			team.RoundWins++
		}
	}

	return round
}

// FinalWinners are the teams with the most round wins. Total votes break a tie; if teams
// are still level they are joint winners.
func (g *Game) FinalWinners() []string {
	winners := make([]string, 0)
	var best *Team
	for _, t := range g.Teams {
		if best == nil || t.RoundWins > best.RoundWins ||
			(t.RoundWins == best.RoundWins && t.TotalVotes > best.TotalVotes) {
			best = t
		}
	}
	if best == nil || (best.RoundWins == 0 && best.TotalVotes == 0) {
		return winners
	}
	for _, t := range g.Teams {
		if t.RoundWins == best.RoundWins && t.TotalVotes == best.TotalVotes {
			winners = append(winners, t.ID)
		}
	}
	return winners
}

// WinningSubmission is the played submission with the best net audience reaction.
// Ties go to the one with more reactions overall, then the one played first.
func (g *Game) WinningSubmission() *Submission {
	var best *Submission
	for _, s := range g.Submissions {
		if s.Status != SUBMISSION_STATUS_PLAYED || s.Played == nil {
			continue
		}
		if best == nil {
			best = s
			continue
		}
		switch {
		case s.NetScore() > best.NetScore():
			best = s
		case s.NetScore() < best.NetScore():
		case len(s.Reactions) > len(best.Reactions):
			best = s
		case len(s.Reactions) == len(best.Reactions) && s.Played.Before(*best.Played):
			best = s
		}
	}
	return best
}

// Reset returns the game to the lobby, keeping players, teams and the submission queue
func (g *Game) Reset() {
	g.GameCode = GAME_CODE_LOBBY
	g.Rounds = make([]*Round, 0)
	g.PinnedTopicID = ""
	for _, t := range g.Teams {
		t.RoundWins = 0
		t.TotalVotes = 0
	}
	for _, s := range g.Submissions {
		if s.Status == SUBMISSION_STATUS_PLAYED {
			s.Status = SUBMISSION_STATUS_APPROVED
			s.Played = nil
			s.Reactions = make(map[string]int)
		}
	}
}
