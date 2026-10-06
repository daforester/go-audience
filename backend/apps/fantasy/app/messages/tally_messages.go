package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

// ReactionTallyMessage shares the live topic's reactions with everyone, at most every
// TallyInterval. Must be called with the game locked.
func ReactionTallyMessage(game *models.Game, submission *models.Submission) {
	roomKey := game.Key
	data := map[string]interface{}{
		"SubmissionID": submission.ID,
		"Ups":          submission.Ups(),
		"Downs":        submission.Downs(),
		"Net":          submission.NetScore(),
	}

	throttled(roomKey+"/reactions", TallyInterval, func() {
		_ = sendToAll(roomKey, "REACTION_TALLY", data)
	})
}

// VoteTallyMessage shares the live vote count with admin and display, and only the number
// of votes cast with everyone else. Must be called with the game locked.
func VoteTallyMessage(game *models.Game, round *models.Round) {
	roomKey := game.Key
	tally := map[string]interface{}{
		"Round":      round.Number,
		"Votes":      round.Tally(),
		"TotalVotes": len(round.Votes),
	}
	count := map[string]interface{}{
		"Round":      round.Number,
		"TotalVotes": len(round.Votes),
	}

	throttled(roomKey+"/votes", TallyInterval, func() {
		_ = sendToClasses(roomKey, "VOTE_TALLY", tally, models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_DISPLAY)
		_ = sendToClasses(roomKey, "VOTE_COUNT", count, models.PLAYER_CLASS_AUDIENCE, models.PLAYER_CLASS_DIRECTOR, models.PLAYER_CLASS_STAFF)
	})
}
