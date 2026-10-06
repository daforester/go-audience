package models

type Round struct {
	Number         int
	SubmissionID   string
	TeamOrder      []string
	PresentingTeam string
	Votes          map[string]string // audience UUID -> team ID
	WinnerTeamIDs  []string
	Decided        bool
}

func NewRound(number int, submissionID string, teamOrder []string) *Round {
	return &Round{
		Number:        number,
		SubmissionID:  submissionID,
		TeamOrder:     teamOrder,
		Votes:         make(map[string]string),
		WinnerTeamIDs: make([]string, 0),
	}
}

func (R *Round) HasTeam(teamID string) bool {
	for _, id := range R.TeamOrder {
		if id == teamID {
			return true
		}
	}
	return false
}

// Vote records one audience member's vote, replacing any earlier one
func (R *Round) Vote(uuid string, teamID string) {
	if R.Votes == nil {
		R.Votes = make(map[string]string)
	}
	R.Votes[uuid] = teamID
}

// Tally counts votes for every team in the round, including teams with none
func (R *Round) Tally() map[string]int {
	tally := make(map[string]int, len(R.TeamOrder))
	for _, id := range R.TeamOrder {
		tally[id] = 0
	}
	for _, teamID := range R.Votes {
		if _, ok := tally[teamID]; ok {
			tally[teamID]++
		}
	}
	return tally
}

// Winners are the teams with the most votes. Tied teams all win. No votes, no winner.
func (R *Round) Winners() []string {
	tally := R.Tally()
	best := 0
	for _, votes := range tally {
		best = max(best, votes)
	}

	winners := make([]string, 0)
	if best == 0 {
		return winners
	}
	for _, id := range R.TeamOrder {
		if tally[id] == best {
			winners = append(winners, id)
		}
	}
	return winners
}
