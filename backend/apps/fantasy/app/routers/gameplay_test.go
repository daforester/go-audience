package routers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	fws "github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
)

// socketLog collects every message a websocket receives
type socketLog struct {
	mux      sync.Mutex
	messages []string
	conn     *fws.Conn
}

func watch(t *testing.T, addr string, p joined) *socketLog {
	t.Helper()
	conn, _, err := dial(addr, "/ws/"+p.RoomKey+"/"+p.UUID+"?token="+p.AccessToken, "http://allowed.example")
	if err != nil {
		t.Fatal(err)
	}
	l := &socketLog{conn: conn}
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			l.mux.Lock()
			l.messages = append(l.messages, string(data))
			l.mux.Unlock()
		}
	}()
	t.Cleanup(func() { _ = conn.Close() })
	return l
}

// waitFor returns the Data of the latest command message matching the predicate
func (l *socketLog) waitFor(t *testing.T, command string, match func(data map[string]interface{}) bool) map[string]interface{} {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		l.mux.Lock()
		for i := len(l.messages) - 1; i >= 0; i-- {
			var msg struct {
				Command string
				Data    map[string]interface{}
			}
			if json.Unmarshal([]byte(l.messages[i]), &msg) == nil && msg.Command == command && (match == nil || match(msg.Data)) {
				l.mux.Unlock()
				return msg.Data
			}
		}
		l.mux.Unlock()
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("no %s message arrived", command)
	return nil
}

func (l *socketLog) saw(substr string) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	for _, m := range l.messages {
		if strings.Contains(m, substr) {
			return true
		}
	}
	return false
}

func gameCode(code float64) func(map[string]interface{}) bool {
	return func(d map[string]interface{}) bool { return d["GameCode"] == code }
}

func (s *testServer) expect(t *testing.T, status int, method string, path string, token string, body interface{}) map[string]interface{} {
	t.Helper()
	code, out, _ := s.do(method, path, token, body, nil)
	if code != status {
		t.Fatalf("%s %s: expected %d, got %d %v", method, path, status, code, out)
	}
	r, _ := out["Response"].(map[string]interface{})
	return r
}

func TestGameplay(t *testing.T) {
	for _, teams := range []int{2, 3} {
		t.Run(fmt.Sprintf("%d teams", teams), func(t *testing.T) {
			playGame(t, teams)
		})
	}
}

func playGame(t *testing.T, teamCount int) {
	s := newTestServer(t, map[string]interface{}{
		"CORS_ALLOW_ORIGINS":          []string{"http://allowed.example"},
		"SUBMISSION_COOLDOWN_SECONDS": 0,
	})
	addr := listen(t, s)

	admin := s.createGame("Admin")
	display := s.joinGame(admin.RoomKey, "Projector")
	staff := s.joinGame(admin.RoomKey, "Staff")
	audience := make([]joined, 4)
	for i := range audience {
		audience[i] = s.joinGame(admin.RoomKey, fmt.Sprintf("Fan %d", i))
	}
	assign := func(p joined, class string, team string) {
		s.expect(t, 200, "POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": p.UUID, "PlayerClass": class, "TeamID": team})
	}
	assign(display, "Display", "")
	assign(staff, "Staff", "")

	// Teams: up to MAX_TEAMS
	teamIDs := make([]string, 0)
	for i := 0; i < teamCount; i++ {
		game := s.expect(t, 200, "POST", "/api/game/team/add", admin.AccessToken, map[string]string{"Name": fmt.Sprintf("Team %c", 'A'+i)})
		teams := game["Game"].(map[string]interface{})["Teams"].([]interface{})
		teamIDs = append(teamIDs, teams[len(teams)-1].(map[string]interface{})["ID"].(string))
	}
	if teamCount == 3 {
		s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/team/add", admin.AccessToken, map[string]string{"Name": "Team D"})
	}
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/team/add", admin.AccessToken, map[string]string{"Name": "team a"})
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/team/add", staff.AccessToken, map[string]string{"Name": "Rogue"})

	directors := make([]joined, teamCount)
	for i := range directors {
		directors[i] = s.joinGame(admin.RoomKey, fmt.Sprintf("Director %d", i))
		assign(directors[i], "Director", teamIDs[i])
	}

	displaySock := watch(t, addr, display)
	staffSock := watch(t, addr, staff)
	fanSock := watch(t, addr, audience[0])
	directorSock := watch(t, addr, directors[0])

	// Lobby: no submissions yet, and you cannot jump straight to a topic
	s.expect(t, fiber.StatusConflict, "PUT", "/api/game/submission/new", audience[0].AccessToken, map[string]string{"Kind": "Complaint", "Text": "Too early"})
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 100})
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 120})

	// Intake opens submissions
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 30})
	fanSock.waitFor(t, "GAME_UPDATE", func(d map[string]interface{}) bool { return d["SubmissionsEnabled"] == true })

	texts := []string{"The lift smells of wet fursuit", "Who stole the con chair's tail?", "Dealers den is on fire (metaphorically)"}
	ids := make([]string, 0)
	for i, text := range texts {
		sub := s.expect(t, 200, "PUT", "/api/game/submission/new", audience[i].AccessToken, map[string]string{"Kind": "Complaint", "Text": "  " + text + "\n"})
		ids = append(ids, sub["Submission"].(map[string]interface{})["ID"].(string))
	}
	s.expect(t, fiber.StatusBadRequest, "PUT", "/api/game/submission/new", audience[0].AccessToken, map[string]string{"Kind": "Gossip", "Text": "x"})
	s.expect(t, fiber.StatusBadRequest, "PUT", "/api/game/submission/new", audience[0].AccessToken, map[string]string{"Kind": "Other", "Text": "   "})
	s.expect(t, fiber.StatusForbidden, "PUT", "/api/game/submission/new", staff.AccessToken, map[string]string{"Kind": "Other", "Text": "staff can't submit"})

	// Staff see new submissions with authors, the audience and display never get the queue
	staffSock.waitFor(t, "NEW_SUBMISSION", func(d map[string]interface{}) bool {
		return d["Submission"].(map[string]interface{})["AuthorName"] == "Fan 1"
	})

	// Staff triage: audience may not, played or unknown ones cannot be voted
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submission/vote/"+ids[1], audience[3].AccessToken, nil)
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submission/status/"+ids[1], directors[0].AccessToken, map[string]string{"Status": "Approved"})
	for _, id := range ids {
		s.expect(t, 200, "POST", "/api/game/submission/status/"+id, staff.AccessToken, map[string]string{"Status": "Approved"})
	}
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/submission/status/"+ids[2], staff.AccessToken, map[string]string{"Status": "Played"})
	sub := s.expect(t, 200, "POST", "/api/game/submission/vote/"+ids[1], staff.AccessToken, nil)
	if sub["Submission"].(map[string]interface{})["StaffVotes"] != 1.0 || sub["Submission"].(map[string]interface{})["MyStaffVote"] != true {
		t.Fatalf("staff vote not recorded: %v", sub)
	}
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submission/vote/"+ids[2], admin.AccessToken, nil) // stars are staff only
	s.expect(t, 200, "POST", "/api/game/submission/vote/"+ids[2], staff.AccessToken, nil)
	s.expect(t, 200, "POST", "/api/game/submission/vote/"+ids[2], staff.AccessToken, nil) // toggled off again
	fanSock.waitFor(t, "SUBMISSION_STATUS", func(d map[string]interface{}) bool { return d["Status"] == "Approved" })

	// Round 1: the most staff-voted approved submission becomes the topic
	game := s.expect(t, 200, "POST", "/api/game/round/start", admin.AccessToken, map[string]string{})
	round := game["Game"].(map[string]interface{})["Round"].(map[string]interface{})
	if round["Topic"].(map[string]interface{})["ID"] != ids[1] {
		t.Fatalf("expected the staff favourite as topic: %v", round)
	}
	reveal := displaySock.waitFor(t, "GAME_UPDATE", gameCode(100))
	topic := reveal["Round"].(map[string]interface{})["Topic"].(map[string]interface{})
	if topic["Text"] != texts[1] || topic["AuthorName"] != nil || topic["AuthorUUID"] != nil {
		t.Fatalf("display topic should be the cleaned text without author: %v", topic)
	}
	if reveal["Submissions"] != nil {
		t.Fatal("display must not receive the submission queue")
	}

	s.expect(t, fiber.StatusConflict, "POST", "/api/game/round/react", audience[0].AccessToken, map[string]int{"Reaction": 1})
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/round/start", admin.AccessToken, map[string]string{})

	// Responses: teams present, audience react
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 110})
	presenting := displaySock.waitFor(t, "GAME_UPDATE", gameCode(110))["Round"].(map[string]interface{})["PresentingTeam"]
	order := round["TeamOrder"].([]interface{})
	if presenting != order[0] {
		t.Fatalf("first team in the order presents first: %v %v", presenting, order)
	}
	s.expect(t, 200, "POST", "/api/game/round/presenting", admin.AccessToken, map[string]string{"TeamID": order[1].(string)})
	directorSock.waitFor(t, "PRESENTING_TEAM", func(d map[string]interface{}) bool { return d["TeamID"] == order[1] })
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/round/presenting", admin.AccessToken, map[string]string{"TeamID": "nope"})

	for i, r := range []int{1, 1, 1, -1} {
		s.expect(t, 200, "POST", "/api/game/round/react", audience[i].AccessToken, map[string]int{"Reaction": r})
	}
	s.expect(t, 200, "POST", "/api/game/round/react", audience[0].AccessToken, map[string]int{"Reaction": 1}) // no double counting
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/round/react", audience[0].AccessToken, map[string]int{"Reaction": 5})
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/round/react", directors[0].AccessToken, map[string]int{"Reaction": 1})
	displaySock.waitFor(t, "REACTION_TALLY", func(d map[string]interface{}) bool { return d["Ups"] == 3.0 && d["Downs"] == 1.0 })
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/round/vote", audience[0].AccessToken, map[string]string{"TeamID": teamIDs[0]})

	// Audience vote: one vote each, changeable, hidden from the audience until decided
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 120})
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/round/react", audience[0].AccessToken, map[string]int{"Reaction": -1})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[0].AccessToken, map[string]string{"TeamID": teamIDs[1]})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[0].AccessToken, map[string]string{"TeamID": teamIDs[0]})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[1].AccessToken, map[string]string{"TeamID": teamIDs[0]})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[2].AccessToken, map[string]string{"TeamID": teamIDs[1]})
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/round/vote", directors[0].AccessToken, map[string]string{"TeamID": teamIDs[0]})
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/round/vote", audience[3].AccessToken, map[string]string{"TeamID": "nope"})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[3].AccessToken, map[string]string{"TeamID": teamIDs[0]})

	displaySock.waitFor(t, "VOTE_TALLY", func(d map[string]interface{}) bool { return d["TotalVotes"] == 4.0 })
	fanSock.waitFor(t, "VOTE_COUNT", func(d map[string]interface{}) bool { return d["TotalVotes"] == 4.0 })
	if fanSock.saw(`"Command":"VOTE_TALLY"`) {
		t.Fatal("audience must not see the live vote tally")
	}
	mine := s.expect(t, 200, "GET", "/api/game", audience[0].AccessToken, nil)["Game"].(map[string]interface{})["Round"].(map[string]interface{})
	if mine["Votes"] != nil || mine["MyVote"] != teamIDs[0] || mine["MyReaction"] != 1.0 {
		t.Fatalf("audience round view wrong: %v", mine)
	}

	s.expect(t, fiber.StatusConflict, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 600})

	// Round result: team 0 wins 3-1
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 130})
	result := fanSock.waitFor(t, "GAME_UPDATE", gameCode(130))["Round"].(map[string]interface{})
	if result["WinnerTeamIDs"].([]interface{})[0] != teamIDs[0] || result["Votes"].(map[string]interface{})[teamIDs[0]] != 3.0 {
		t.Fatalf("round result wrong: %v", result)
	}
	s.expect(t, fiber.StatusConflict, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 120})

	// Round 2 straight from the result, with the admin picking the topic
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/round/start", admin.AccessToken, map[string]string{"SubmissionID": ids[1]}) // already played
	s.expect(t, 200, "POST", "/api/game/round/start", admin.AccessToken, map[string]string{"SubmissionID": ids[0]})
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 110})
	s.expect(t, 200, "POST", "/api/game/round/react", audience[0].AccessToken, map[string]int{"Reaction": -1})
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 120})
	s.expect(t, 200, "POST", "/api/game/round/vote", audience[0].AccessToken, map[string]string{"TeamID": teamIDs[1]})
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 130})

	// Final tally: team 0 and 1 each won a round, team 0 has more total votes (3 to 2)
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 600})
	final := displaySock.waitFor(t, "GAME_UPDATE", gameCode(600))["Final"].(map[string]interface{})
	if final["WinnerTeamIDs"].([]interface{})[0] != teamIDs[0] || len(final["WinnerTeamIDs"].([]interface{})) != 1 {
		t.Fatalf("final winners wrong: %v", final)
	}
	best := final["WinningSubmission"].(map[string]interface{})
	if best["ID"] != ids[1] || best["Net"] != 2.0 || best["AuthorName"] != nil {
		t.Fatalf("winning submission wrong: %v", best)
	}

	// The admin reveals who wrote it
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submitter-names", staff.AccessToken, map[string]bool{"Enabled": true})
	s.expect(t, 200, "POST", "/api/game/submitter-names", admin.AccessToken, map[string]bool{"Enabled": true})
	named := displaySock.waitFor(t, "GAME_UPDATE", func(d map[string]interface{}) bool { return d["ShowSubmitterNames"] == true })
	if named["Final"].(map[string]interface{})["WinningSubmission"].(map[string]interface{})["AuthorName"] != "Fan 1" {
		t.Fatalf("author should show once names are on: %v", named["Final"])
	}

	for _, sock := range []*socketLog{fanSock, directorSock, displaySock} {
		if sock.saw("NEW_SUBMISSION") || sock.saw("SUBMISSION_UPDATED") {
			t.Fatal("only admin and staff receive queue messages")
		}
	}
	if fanSock.saw(`"AuthorUUID"`) || displaySock.saw(`"AuthorUUID"`) {
		t.Fatal("author UUIDs never leave admin and staff")
	}

	// Reset back to the lobby keeps players and teams but clears scores
	reset := s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 20})["Game"].(map[string]interface{})
	if reset["Round"] != nil || reset["Teams"].([]interface{})[0].(map[string]interface{})["RoundWins"] != 0.0 {
		t.Fatalf("reset should clear rounds and scores: %v", reset)
	}

	// Removing a team sends its directors back to the audience
	s.expect(t, 200, "POST", "/api/game/team/remove/"+teamIDs[0], admin.AccessToken, nil)
	directorSock.waitFor(t, "PLAYER_UPDATE", func(d map[string]interface{}) bool { return d["Class"] == "Audience" })
}

func TestSubmissionCooldown(t *testing.T) {
	s := newTestServer(t, map[string]interface{}{"SUBMISSION_COOLDOWN_SECONDS": 60})
	admin := s.createGame("Admin")
	fan := s.joinGame(admin.RoomKey, "Fan")
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 30})

	s.expect(t, 200, "PUT", "/api/game/submission/new", fan.AccessToken, map[string]string{"Kind": "Question", "Text": "One"})
	s.expect(t, fiber.StatusTooManyRequests, "PUT", "/api/game/submission/new", fan.AccessToken, map[string]string{"Kind": "Question", "Text": "Two"})

	long := strings.Repeat("a", 1000)
	other := s.joinGame(admin.RoomKey, "Other")
	sub := s.expect(t, 200, "PUT", "/api/game/submission/new", other.AccessToken, map[string]string{"Kind": "Question", "Text": long})
	if len(sub["Submission"].(map[string]interface{})["Text"].(string)) != 280 {
		t.Fatal("submissions are cut to MAX_SUBMISSION_LENGTH")
	}

	s.expect(t, 200, "POST", "/api/game/submissions", admin.AccessToken, map[string]bool{"Enabled": false})
	late := s.joinGame(admin.RoomKey, "Late")
	s.expect(t, fiber.StatusConflict, "PUT", "/api/game/submission/new", late.AccessToken, map[string]string{"Kind": "Question", "Text": "Closed?"})
}

func TestAdminTopPicksVetoAndOverride(t *testing.T) {
	s := newTestServer(t, map[string]interface{}{
		"CORS_ALLOW_ORIGINS":          []string{"http://allowed.example"},
		"SUBMISSION_COOLDOWN_SECONDS": 0,
	})
	addr := listen(t, s)

	admin := s.createGame("Admin")
	staff := s.joinGame(admin.RoomKey, "Staff")
	fan := s.joinGame(admin.RoomKey, "Fan")
	s.expect(t, 200, "POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": staff.UUID, "PlayerClass": "Staff"})
	for i, name := range []string{"Team A", "Team B"} {
		game := s.expect(t, 200, "POST", "/api/game/team/add", admin.AccessToken, map[string]string{"Name": name})
		teams := game["Game"].(map[string]interface{})["Teams"].([]interface{})
		director := s.joinGame(admin.RoomKey, fmt.Sprintf("Director %d", i))
		teamID := teams[i].(map[string]interface{})["ID"].(string)
		s.expect(t, 200, "POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": director.UUID, "PlayerClass": "Director", "TeamID": teamID})
	}
	staffSock := watch(t, addr, staff)
	s.expect(t, 200, "POST", "/api/game/code", admin.AccessToken, map[string]int{"Code": 30})

	ids := make([]string, 0)
	for i := 0; i < 4; i++ {
		sub := s.expect(t, 200, "PUT", "/api/game/submission/new", fan.AccessToken, map[string]string{"Kind": "Question", "Text": fmt.Sprintf("Question %d", i)})
		ids = append(ids, sub["Submission"].(map[string]interface{})["ID"].(string))
	}
	for _, id := range ids[:3] {
		s.expect(t, 200, "POST", "/api/game/submission/status/"+id, staff.AccessToken, map[string]string{"Status": "Approved"})
	}
	s.expect(t, 200, "POST", "/api/game/submission/vote/"+ids[2], staff.AccessToken, nil)

	// The admin sees the staff's top picks in play order
	game := s.expect(t, 200, "GET", "/api/game", admin.AccessToken, nil)["Game"].(map[string]interface{})
	picks := game["TopPickIDs"].([]interface{})
	if len(picks) != 3 || picks[0] != ids[2] || picks[1] != ids[0] || game["NextTopicID"] != ids[2] {
		t.Fatalf("top picks wrong: %v next %v", picks, game["NextTopicID"])
	}

	// Override: pin a submission still awaiting review, it becomes the next topic
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submission/pin", staff.AccessToken, map[string]string{"SubmissionID": ids[3]})
	s.expect(t, 200, "POST", "/api/game/submission/pin", admin.AccessToken, map[string]string{"SubmissionID": ids[3]})
	staffSock.waitFor(t, "QUEUE_UPDATED", func(d map[string]interface{}) bool {
		return d["PinnedTopicID"] == ids[3] && d["NextTopicID"] == ids[3]
	})

	// Veto: the admin rejects the staff favourite; staff cannot restore it, the admin can
	s.expect(t, 200, "POST", "/api/game/submission/status/"+ids[2], admin.AccessToken, map[string]string{"Status": "Rejected"})
	s.expect(t, fiber.StatusForbidden, "POST", "/api/game/submission/status/"+ids[2], staff.AccessToken, map[string]string{"Status": "Approved"})
	vetoed := staffSock.waitFor(t, "SUBMISSION_UPDATED", func(d map[string]interface{}) bool {
		sub := d["Submission"].(map[string]interface{})
		return sub["ID"] == ids[2] && sub["Vetoed"] == true
	})
	if picks := vetoed["TopPickIDs"].([]interface{}); len(picks) != 2 || picks[0] != ids[0] {
		t.Fatalf("vetoed submission should leave the top picks: %v", picks)
	}

	// Vetoing the pinned submission removes the pin
	s.expect(t, 200, "POST", "/api/game/submission/status/"+ids[3], admin.AccessToken, map[string]string{"Status": "Rejected"})
	game = s.expect(t, 200, "GET", "/api/game", admin.AccessToken, nil)["Game"].(map[string]interface{})
	if game["PinnedTopicID"] != nil || game["NextTopicID"] != ids[0] {
		t.Fatalf("veto should clear the pin: pinned %v next %v", game["PinnedTopicID"], game["NextTopicID"])
	}
	s.expect(t, fiber.StatusBadRequest, "POST", "/api/game/submission/pin", admin.AccessToken, map[string]string{"SubmissionID": ids[3]})

	// The admin lifts the veto and pins it again; revealing plays the pin, then clears it
	s.expect(t, 200, "POST", "/api/game/submission/status/"+ids[2], admin.AccessToken, map[string]string{"Status": "Approved"})
	s.expect(t, 200, "POST", "/api/game/submission/status/"+ids[2], staff.AccessToken, map[string]string{"Status": "Pending"})
	s.expect(t, 200, "POST", "/api/game/submission/pin", admin.AccessToken, map[string]string{"SubmissionID": ids[2]})
	started := s.expect(t, 200, "POST", "/api/game/round/start", admin.AccessToken, map[string]string{})["Game"].(map[string]interface{})
	topic := started["Round"].(map[string]interface{})["Topic"].(map[string]interface{})
	if topic["ID"] != ids[2] || started["PinnedTopicID"] != nil {
		t.Fatalf("pinned submission should be played and unpinned: %v", started["Round"])
	}

	// The audience never sees the queue order
	fanView := s.expect(t, 200, "GET", "/api/game", fan.AccessToken, nil)["Game"].(map[string]interface{})
	if fanView["TopPickIDs"] != nil || fanView["PinnedTopicID"] != nil {
		t.Fatal("queue order is for admin and staff only")
	}
}
