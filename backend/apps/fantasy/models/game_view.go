package models

import "time"

// Views are what a player is allowed to see of a game. They are always built from copies,
// so they can be marshalled after the game lock has been released.

type PlayerSummary struct {
	UUID string
	Name string
}

type TeamView struct {
	ID         string
	Name       string
	RoundWins  int
	TotalVotes int
	Members    []PlayerSummary
}

type SubmissionView struct {
	ID          string
	Kind        SubmissionKind
	Text        string
	AuthorName  string `json:",omitempty"`
	AuthorUUID  string `json:",omitempty"`
	Status      SubmissionStatus
	StaffVotes  int
	MyStaffVote bool `json:",omitempty"`
	Ups         int
	Downs       int
	Net         int
	Created     time.Time
}

type RoundView struct {
	Number         int
	Topic          *SubmissionView
	TeamOrder      []string
	PresentingTeam string
	Votes          map[string]int `json:",omitempty"`
	TotalVotes     int
	WinnerTeamIDs  []string
	Decided        bool
	MyVote         string `json:",omitempty"`
	MyReaction     int    `json:",omitempty"`
}

type FinalView struct {
	WinnerTeamIDs     []string
	WinningSubmission *SubmissionView
}

type GameView struct {
	Key                string
	GameCode           int
	SubmissionsEnabled bool
	ShowSubmitterNames bool
	Volume             int
	MaxTeams           int `json:",omitempty"`
	SubmissionKinds    []SubmissionKind
	Players            []*Player `json:",omitempty"`
	Teams              []TeamView
	Submissions        []SubmissionView `json:",omitempty"`
	MySubmissions      []SubmissionView `json:",omitempty"`
	NextTopicID        string           `json:",omitempty"`
	Round              *RoundView       `json:",omitempty"`
	Final              *FinalView       `json:",omitempty"`
}

func canModerate(class PlayerClass) bool {
	return class == PLAYER_CLASS_ADMIN || class == PLAYER_CLASS_STAFF
}

func canSeeLiveVotes(class PlayerClass) bool {
	return class == PLAYER_CLASS_ADMIN || class == PLAYER_CLASS_DISPLAY
}

// ViewFor builds the view a player may see. A player without a UUID gets the shared view
// for their class, without personal fields such as their own votes.
func (g *Game) ViewFor(player *Player) *GameView {
	view := &GameView{
		Key:                g.Key,
		GameCode:           g.GameCode,
		SubmissionsEnabled: g.SubmissionsEnabled,
		ShowSubmitterNames: g.ShowSubmitterNames,
		Volume:             g.Volume,
		SubmissionKinds:    SubmissionKinds,
		Teams:              g.TeamViews(),
	}

	if player.Class == PLAYER_CLASS_ADMIN || player.Class == PLAYER_CLASS_DISPLAY {
		view.Players = g.CopyPlayers()
	}

	if canModerate(player.Class) {
		view.Submissions = make([]SubmissionView, 0, len(g.Submissions))
		for _, s := range g.Submissions {
			view.Submissions = append(view.Submissions, *g.SubmissionViewFor(s, player))
		}
		if next := g.NextTopic(); next != nil {
			view.NextTopicID = next.ID
		}
	}

	if player.UUID != "" && player.Class == PLAYER_CLASS_AUDIENCE {
		view.MySubmissions = make([]SubmissionView, 0)
		for _, s := range g.Submissions {
			if s.AuthorUUID == player.UUID {
				view.MySubmissions = append(view.MySubmissions, *g.SubmissionViewFor(s, player))
			}
		}
	}

	view.Round = g.RoundViewFor(player)

	if g.GameCode == GAME_CODE_FINAL_TALLY {
		view.Final = &FinalView{
			WinnerTeamIDs: g.FinalWinners(),
		}
		if s := g.WinningSubmission(); s != nil {
			view.Final.WinningSubmission = g.SubmissionViewFor(s, player)
		}
	}

	return view
}

func (g *Game) TeamViews() []TeamView {
	teams := make([]TeamView, 0, len(g.Teams))
	for _, t := range g.Teams {
		members := make([]PlayerSummary, 0)
		for _, p := range g.GetTeamDirectors(t.ID) {
			members = append(members, PlayerSummary{p.UUID, p.Name})
		}
		teams = append(teams, TeamView{
			ID:         t.ID,
			Name:       t.Name,
			RoundWins:  t.RoundWins,
			TotalVotes: t.TotalVotes,
			Members:    members,
		})
	}
	return teams
}

// SubmissionViewFor applies the submitter names setting: admin and staff always see who
// wrote a submission, everyone else only when the admin has turned names on
func (g *Game) SubmissionViewFor(s *Submission, player *Player) *SubmissionView {
	view := &SubmissionView{
		ID:         s.ID,
		Kind:       s.Kind,
		Text:       s.Text,
		Status:     s.Status,
		StaffVotes: s.StaffVoteCount(),
		Ups:        s.Ups(),
		Downs:      s.Downs(),
		Net:        s.NetScore(),
		Created:    s.Created,
	}

	if canModerate(player.Class) {
		view.AuthorName = s.AuthorName
		view.AuthorUUID = s.AuthorUUID
		view.MyStaffVote = player.UUID != "" && s.StaffVotes[player.UUID]
	} else if g.ShowSubmitterNames || (player.UUID != "" && s.AuthorUUID == player.UUID) {
		view.AuthorName = s.AuthorName
	}

	return view
}

// RoundViewFor hides the live vote count from everyone except the admin and display
// until the round is decided, so nobody votes with the crowd
func (g *Game) RoundViewFor(player *Player) *RoundView {
	round := g.GetCurrentRound()
	if round == nil {
		return nil
	}

	view := &RoundView{
		Number:         round.Number,
		TeamOrder:      append([]string{}, round.TeamOrder...),
		PresentingTeam: round.PresentingTeam,
		TotalVotes:     len(round.Votes),
		WinnerTeamIDs:  append([]string{}, round.WinnerTeamIDs...),
		Decided:        round.Decided,
	}

	if s := g.GetSubmission(round.SubmissionID); s != nil {
		view.Topic = g.SubmissionViewFor(s, player)
		if player.UUID != "" {
			view.MyReaction = s.Reactions[player.UUID]
		}
	}

	if round.Decided || canSeeLiveVotes(player.Class) {
		view.Votes = round.Tally()
	}

	if player.UUID != "" {
		view.MyVote = round.Votes[player.UUID]
	}

	return view
}
