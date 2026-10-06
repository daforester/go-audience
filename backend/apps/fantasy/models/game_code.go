package models

const (
	GAME_CODE_LOBBY         = 20
	GAME_CODE_INTAKE        = 30
	GAME_CODE_TOPIC_REVEAL  = 100
	GAME_CODE_RESPONSES     = 110
	GAME_CODE_AUDIENCE_VOTE = 120
	GAME_CODE_ROUND_RESULT  = 130
	GAME_CODE_FINAL_TALLY   = 600
)

// GameCodeTransitions lists the codes each code may move to. Entering
// GAME_CODE_TOPIC_REVEAL starts a round, so it is only reached through StartRound.
// Any code may return to GAME_CODE_LOBBY, which resets the game.
var GameCodeTransitions = map[int][]int{
	GAME_CODE_LOBBY:         {GAME_CODE_INTAKE},
	GAME_CODE_INTAKE:        {GAME_CODE_TOPIC_REVEAL, GAME_CODE_FINAL_TALLY},
	GAME_CODE_TOPIC_REVEAL:  {GAME_CODE_RESPONSES},
	GAME_CODE_RESPONSES:     {GAME_CODE_AUDIENCE_VOTE},
	GAME_CODE_AUDIENCE_VOTE: {GAME_CODE_ROUND_RESULT},
	GAME_CODE_ROUND_RESULT:  {GAME_CODE_INTAKE, GAME_CODE_TOPIC_REVEAL, GAME_CODE_FINAL_TALLY},
	GAME_CODE_FINAL_TALLY:   {},
}

func CanTransition(from int, to int) bool {
	if to == GAME_CODE_LOBBY {
		return true
	}
	for _, code := range GameCodeTransitions[from] {
		if code == to {
			return true
		}
	}
	return false
}
