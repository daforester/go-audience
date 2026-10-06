package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

// moderators are sent their own view of a submission, which includes their staff vote
func sendToModerators(game *models.Game, command string, submission *models.Submission) error {
	for _, player := range game.Players {
		if player.Class != models.PLAYER_CLASS_ADMIN && player.Class != models.PLAYER_CLASS_STAFF {
			continue
		}
		next, pinned, picks := game.QueueOrder()
		err := sendToPlayer(game.Key, player.UUID, command, struct {
			Submission    *models.SubmissionView
			NextTopicID   string
			PinnedTopicID string
			TopPickIDs    []string
		}{
			game.SubmissionViewFor(submission, player),
			next,
			pinned,
			picks,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// QueueUpdatedMessage tells admin and staff the queue order changed, e.g. a pin. Must be
// called with the game locked.
func QueueUpdatedMessage(game *models.Game) error {
	next, pinned, picks := game.QueueOrder()
	return sendToClasses(game.Key, "QUEUE_UPDATED", map[string]interface{}{
		"NextTopicID":   next,
		"PinnedTopicID": pinned,
		"TopPickIDs":    picks,
	}, models.PLAYER_CLASS_ADMIN, models.PLAYER_CLASS_STAFF)
}

// NewSubmissionMessage tells admin and staff about a new submission. Must be called with
// the game locked.
func NewSubmissionMessage(game *models.Game, submission *models.Submission) error {
	SubmissionCountMessage(game)
	return sendToModerators(game, "NEW_SUBMISSION", submission)
}

// SubmissionCountMessage tells the display how many submissions have arrived, at most every
// TallyInterval. Must be called with the game locked.
func SubmissionCountMessage(game *models.Game) {
	roomKey := game.Key
	data := map[string]interface{}{
		"SubmissionCount": len(game.Submissions),
	}

	throttled(roomKey+"/submissions", TallyInterval, func() {
		_ = sendToClasses(roomKey, "SUBMISSION_COUNT", data, models.PLAYER_CLASS_DISPLAY, models.PLAYER_CLASS_ADMIN)
	})
}

// SubmissionUpdatedMessage tells admin and staff a submission changed, and lets the author
// know its status. Must be called with the game locked.
func SubmissionUpdatedMessage(game *models.Game, submission *models.Submission) error {
	err := sendToModerators(game, "SUBMISSION_UPDATED", submission)
	if err != nil {
		return err
	}

	return sendToPlayer(game.Key, submission.AuthorUUID, "SUBMISSION_STATUS", map[string]interface{}{
		"ID":     submission.ID,
		"Status": submission.Status,
	})
}
