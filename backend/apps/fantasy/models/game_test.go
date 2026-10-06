package models

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func testGame(teams int) *Game {
	g := NewGame()
	g.Key = "TEST1"
	g.AddPlayer(&Player{UUID: "admin", Name: "Admin", Class: PLAYER_CLASS_ADMIN})
	g.AddPlayer(&Player{UUID: "display", Name: "Display", Class: PLAYER_CLASS_DISPLAY})
	g.AddPlayer(&Player{UUID: "staff1", Name: "Staff One", Class: PLAYER_CLASS_STAFF})
	g.AddPlayer(&Player{UUID: "staff2", Name: "Staff Two", Class: PLAYER_CLASS_STAFF})
	for i := 0; i < teams; i++ {
		id := fmt.Sprintf("team%d", i+1)
		g.AddTeam(&Team{ID: id, Name: fmt.Sprintf("Team %d", i+1)})
		g.AddPlayer(&Player{UUID: "director-" + id, Name: "Director " + id, Class: PLAYER_CLASS_DIRECTOR, TeamID: id})
	}
	for i := 0; i < 10; i++ {
		g.AddPlayer(&Player{UUID: fmt.Sprintf("aud%d", i), Name: fmt.Sprintf("Audience %d", i), Class: PLAYER_CLASS_AUDIENCE})
	}
	return g
}

func submit(g *Game, id string, author string, created time.Time) *Submission {
	s := NewSubmission(id, g.GetPlayerByUUID(author), SUBMISSION_KIND_COMPLAINT, "The "+id+" is broken")
	s.Created = created
	g.AddSubmission(s)
	return s
}

func TestNextTopic(t *testing.T) {
	g := testGame(2)
	now := time.Now()
	a := submit(g, "a", "aud0", now)
	b := submit(g, "b", "aud1", now.Add(time.Second))
	c := submit(g, "c", "aud2", now.Add(2*time.Second))

	if g.NextTopic() != nil {
		t.Fatal("pending submissions must not be picked")
	}

	a.Status, b.Status, c.Status = SUBMISSION_STATUS_APPROVED, SUBMISSION_STATUS_APPROVED, SUBMISSION_STATUS_REJECTED
	if g.NextTopic() != a {
		t.Fatal("oldest wins a tie")
	}

	b.ToggleStaffVote("staff1")
	c.ToggleStaffVote("staff1")
	c.ToggleStaffVote("staff2")
	if g.NextTopic() != b {
		t.Fatal("most staff votes among approved should win, rejected never")
	}

	b.ToggleStaffVote("staff1")
	if b.StaffVoteCount() != 0 || g.NextTopic() != a {
		t.Fatal("toggling a staff vote off should remove it")
	}
}

func TestVotesAndReactionsAreOnePerPerson(t *testing.T) {
	s := NewSubmission("s", &Player{UUID: "aud0"}, SUBMISSION_KIND_OTHER, "x")
	s.SetReaction("aud1", 1)
	s.SetReaction("aud1", 5)
	s.SetReaction("aud2", -1)
	s.SetReaction("aud3", 1)
	s.SetReaction("aud3", 0)
	if s.Ups() != 1 || s.Downs() != 1 || s.NetScore() != 0 {
		t.Fatalf("reactions: ups %d downs %d", s.Ups(), s.Downs())
	}

	r := NewRound(1, "s", []string{"team1", "team2"})
	r.Vote("aud1", "team1")
	r.Vote("aud1", "team1")
	r.Vote("aud2", "team1")
	r.Vote("aud2", "team2")
	if !reflect.DeepEqual(r.Tally(), map[string]int{"team1": 1, "team2": 1}) {
		t.Fatalf("tally: %v", r.Tally())
	}
}

func TestRoundWinners(t *testing.T) {
	r := NewRound(1, "s", []string{"team1", "team2", "team3"})
	if len(r.Winners()) != 0 {
		t.Fatal("no votes, no winner")
	}
	r.Vote("a", "team2")
	if !reflect.DeepEqual(r.Winners(), []string{"team2"}) {
		t.Fatal(r.Winners())
	}
	r.Vote("b", "team3")
	if !reflect.DeepEqual(r.Winners(), []string{"team2", "team3"}) {
		t.Fatal("tied teams all win", r.Winners())
	}
	r.Vote("c", "unknown-team")
	if r.Tally()["unknown-team"] != 0 {
		t.Fatal("votes for teams not in the round are ignored")
	}
}

func TestStartRoundRotatesTeams(t *testing.T) {
	g := testGame(3)
	now := time.Now()
	for i := 0; i < 4; i++ {
		s := submit(g, fmt.Sprintf("s%d", i), "aud0", now.Add(time.Duration(i)*time.Second))
		s.Status = SUBMISSION_STATUS_APPROVED
	}

	expected := [][]string{
		{"team1", "team2", "team3"},
		{"team2", "team3", "team1"},
		{"team3", "team1", "team2"},
		{"team1", "team2", "team3"},
	}
	for i, order := range expected {
		topic := g.NextTopic()
		r := g.StartRound(topic)
		if topic.Status != SUBMISSION_STATUS_PLAYED || topic.Played == nil {
			t.Fatal("started topic should be marked played")
		}
		if r.Number != i+1 || !reflect.DeepEqual(r.TeamOrder, order) {
			t.Fatalf("round %d order %v, expected %v", r.Number, r.TeamOrder, order)
		}
	}
}

func TestDecideRoundCountsOnce(t *testing.T) {
	g := testGame(2)
	s := submit(g, "s", "aud0", time.Now())
	s.Status = SUBMISSION_STATUS_APPROVED
	r := g.StartRound(s)
	r.Vote("aud1", "team1")
	r.Vote("aud2", "team1")
	r.Vote("aud3", "team2")

	g.DecideRound()
	g.DecideRound()

	if g.GetTeam("team1").RoundWins != 1 || g.GetTeam("team1").TotalVotes != 2 || g.GetTeam("team2").TotalVotes != 1 {
		t.Fatalf("scores: %+v %+v", g.GetTeam("team1"), g.GetTeam("team2"))
	}
}

func TestFinalWinners(t *testing.T) {
	g := testGame(3)
	if len(g.FinalWinners()) != 0 {
		t.Fatal("nobody wins a game with no votes")
	}

	g.GetTeam("team1").RoundWins, g.GetTeam("team1").TotalVotes = 2, 10
	g.GetTeam("team2").RoundWins, g.GetTeam("team2").TotalVotes = 2, 12
	g.GetTeam("team3").RoundWins, g.GetTeam("team3").TotalVotes = 1, 30
	if !reflect.DeepEqual(g.FinalWinners(), []string{"team2"}) {
		t.Fatal("round wins first, then total votes", g.FinalWinners())
	}

	g.GetTeam("team1").TotalVotes = 12
	if !reflect.DeepEqual(g.FinalWinners(), []string{"team1", "team2"}) {
		t.Fatal("still level means joint winners", g.FinalWinners())
	}
}

func TestWinningSubmission(t *testing.T) {
	g := testGame(2)
	now := time.Now()
	a := submit(g, "a", "aud0", now)
	b := submit(g, "b", "aud1", now)
	c := submit(g, "c", "aud2", now)
	unplayed := submit(g, "d", "aud3", now)
	unplayed.SetReaction("aud9", 1)

	for i, s := range []*Submission{a, b, c} {
		played := now.Add(time.Duration(i) * time.Minute)
		s.Status, s.Played = SUBMISSION_STATUS_PLAYED, &played
	}

	a.SetReaction("aud1", 1)
	b.SetReaction("aud1", 1)
	b.SetReaction("aud2", 1)
	b.SetReaction("aud3", -1)
	c.SetReaction("aud4", -1)

	// a: +1 from 1, b: +1 from 3, c: -1
	if g.WinningSubmission() != b {
		t.Fatal("equal net score goes to more reactions")
	}

	b.SetReaction("aud3", 0)
	b.SetReaction("aud2", 0)
	if g.WinningSubmission() != a {
		t.Fatal("fully tied goes to the one played first")
	}

	c.SetReaction("aud4", 1)
	c.SetReaction("aud5", 1)
	if g.WinningSubmission() != c {
		t.Fatal("best net score wins")
	}
}

func TestViewFiltering(t *testing.T) {
	g := testGame(2)
	s := submit(g, "s", "aud0", time.Now())
	s.Status = SUBMISSION_STATUS_APPROVED
	s.ToggleStaffVote("staff1")
	r := g.StartRound(s)
	g.GameCode = GAME_CODE_AUDIENCE_VOTE
	r.Vote("aud1", "team2")

	admin := g.ViewFor(g.GetPlayerByUUID("admin"))
	staff := g.ViewFor(g.GetPlayerByUUID("staff1"))
	display := g.ViewFor(&Player{Class: PLAYER_CLASS_DISPLAY})
	audience := g.ViewFor(g.GetPlayerByUUID("aud1"))
	author := g.ViewFor(g.GetPlayerByUUID("aud0"))
	director := g.ViewFor(g.GetPlayerByUUID("director-team1"))

	if len(admin.Submissions) != 1 || len(staff.Submissions) != 1 || staff.Submissions[0].AuthorName != "Audience 0" || !staff.Submissions[0].MyStaffVote {
		t.Fatal("admin and staff see the queue with authors")
	}
	if len(display.Submissions) != 0 || len(audience.Submissions) != 0 || len(director.Submissions) != 0 {
		t.Fatal("only admin and staff see the queue")
	}
	if len(audience.Players) != 0 || len(director.Players) != 0 || len(admin.Players) == 0 {
		t.Fatal("only admin and display see the player list")
	}

	if display.Round.Topic.AuthorName != "" || audience.Round.Topic.AuthorName != "" {
		t.Fatal("submitter names are hidden by default")
	}
	if author.Round.Topic.AuthorName != "Audience 0" || len(author.MySubmissions) != 1 {
		t.Fatal("authors see their own submissions")
	}

	if display.Round.Votes == nil || admin.Round.Votes["team2"] != 1 {
		t.Fatal("admin and display see live votes")
	}
	if audience.Round.Votes != nil || director.Round.Votes != nil {
		t.Fatal("live votes are hidden from audience and directors")
	}
	if audience.Round.MyVote != "team2" {
		t.Fatal("audience see their own vote")
	}

	g.ShowSubmitterNames = true
	if g.ViewFor(&Player{Class: PLAYER_CLASS_DISPLAY}).Round.Topic.AuthorName != "Audience 0" {
		t.Fatal("submitter names shown when enabled")
	}
	if g.ViewFor(&Player{Class: PLAYER_CLASS_DISPLAY}).Round.Topic.AuthorUUID != "" {
		t.Fatal("author UUIDs are only for admin and staff")
	}

	g.DecideRound()
	if g.ViewFor(g.GetPlayerByUUID("aud1")).Round.Votes["team2"] != 1 {
		t.Fatal("votes are visible to all once decided")
	}
}

func TestTransitions(t *testing.T) {
	allowed := [][2]int{
		{GAME_CODE_LOBBY, GAME_CODE_INTAKE},
		{GAME_CODE_INTAKE, GAME_CODE_TOPIC_REVEAL},
		{GAME_CODE_TOPIC_REVEAL, GAME_CODE_RESPONSES},
		{GAME_CODE_RESPONSES, GAME_CODE_AUDIENCE_VOTE},
		{GAME_CODE_AUDIENCE_VOTE, GAME_CODE_ROUND_RESULT},
		{GAME_CODE_ROUND_RESULT, GAME_CODE_TOPIC_REVEAL},
		{GAME_CODE_ROUND_RESULT, GAME_CODE_FINAL_TALLY},
		{GAME_CODE_FINAL_TALLY, GAME_CODE_LOBBY},
		{GAME_CODE_RESPONSES, GAME_CODE_LOBBY},
	}
	for _, a := range allowed {
		if !CanTransition(a[0], a[1]) {
			t.Errorf("%d -> %d should be allowed", a[0], a[1])
		}
	}
	denied := [][2]int{
		{GAME_CODE_LOBBY, GAME_CODE_TOPIC_REVEAL},
		{GAME_CODE_TOPIC_REVEAL, GAME_CODE_AUDIENCE_VOTE},
		{GAME_CODE_RESPONSES, GAME_CODE_ROUND_RESULT},
		{GAME_CODE_AUDIENCE_VOTE, GAME_CODE_FINAL_TALLY},
		{GAME_CODE_FINAL_TALLY, GAME_CODE_INTAKE},
		{GAME_CODE_INTAKE, 12345},
	}
	for _, d := range denied {
		if CanTransition(d[0], d[1]) {
			t.Errorf("%d -> %d should be denied", d[0], d[1])
		}
	}
}

func TestReadyForRound(t *testing.T) {
	g := testGame(1)
	if g.ReadyForRound() == nil {
		t.Fatal("one team is not enough")
	}
	g.AddTeam(&Team{ID: "empty", Name: "Empty"})
	if g.ReadyForRound() == nil {
		t.Fatal("every team needs a director")
	}
	g.RemoveTeam("empty")
	g.AddTeam(&Team{ID: "team2", Name: "Team 2"})
	g.AddPlayer(&Player{UUID: "d2", Class: PLAYER_CLASS_DIRECTOR, TeamID: "team2"})
	if err := g.ReadyForRound(); err != nil {
		t.Fatal(err)
	}

	g.RemoveTeam("team2")
	if p := g.GetPlayerByUUID("d2"); p.Class != PLAYER_CLASS_AUDIENCE || p.TeamID != "" {
		t.Fatal("removing a team returns its directors to the audience")
	}
}

// A whole game with N teams: rounds are played, scores tallied, winners revealed
func TestFullGame(t *testing.T) {
	for _, teams := range []int{2, 3} {
		t.Run(fmt.Sprintf("%d teams", teams), func(t *testing.T) {
			g := testGame(teams)
			now := time.Now()
			for i := 0; i < 5; i++ {
				s := submit(g, fmt.Sprintf("s%d", i), fmt.Sprintf("aud%d", i), now.Add(time.Duration(i)*time.Second))
				s.Status = SUBMISSION_STATUS_APPROVED
				for j := 0; j < 5-i; j++ {
					s.ToggleStaffVote(fmt.Sprintf("staff%d", j))
				}
			}

			// Round n is won by team (n % teams); topic s1 gets the most love
			for n := 0; n < 3; n++ {
				if err := g.ReadyForRound(); err != nil {
					t.Fatal(err)
				}
				topic := g.NextTopic()
				if topic.ID != fmt.Sprintf("s%d", n) {
					t.Fatalf("round %d expected topic s%d, got %s", n+1, n, topic.ID)
				}
				r := g.StartRound(topic)
				for a := 0; a < 10; a++ {
					reaction := -1
					if topic.ID == "s1" || a%3 == 0 {
						reaction = 1
					}
					topic.SetReaction(fmt.Sprintf("aud%d", a), reaction)
				}
				winner := fmt.Sprintf("team%d", n%teams+1)
				for a := 0; a < 10; a++ {
					team := winner
					if a < 3 {
						team = r.TeamOrder[(a+1)%teams]
						if team == winner {
							team = r.TeamOrder[(a+2)%teams]
						}
					}
					r.Vote(fmt.Sprintf("aud%d", a), team)
				}
				g.DecideRound()
				if !reflect.DeepEqual(r.WinnerTeamIDs, []string{winner}) {
					t.Fatalf("round %d winners %v, expected %s", n+1, r.WinnerTeamIDs, winner)
				}
			}

			g.GameCode = GAME_CODE_FINAL_TALLY
			view := g.ViewFor(&Player{Class: PLAYER_CLASS_DISPLAY})
			if view.Final == nil || view.Final.WinningSubmission.ID != "s1" {
				t.Fatalf("winning submission should be s1: %+v", view.Final)
			}
			// Two teams: team1 wins rounds 1 and 3. Three teams: one round each, equal votes.
			expected := []string{"team1"}
			if teams == 3 {
				expected = []string{"team1", "team2", "team3"}
			}
			if !reflect.DeepEqual(view.Final.WinnerTeamIDs, expected) {
				t.Fatalf("final winners %v", view.Final.WinnerTeamIDs)
			}
			if len(view.Teams) != teams {
				t.Fatalf("expected %d teams in the view", teams)
			}

			g.Reset()
			if len(g.Rounds) != 0 || g.GetTeam("team1").RoundWins != 0 || g.GetSubmission("s0").Status != SUBMISSION_STATUS_APPROVED {
				t.Fatal("reset clears rounds and scores, played topics go back to the queue")
			}
		})
	}
}
