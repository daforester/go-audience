package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type ChangeSubmissionStatusProcessorInterface interface {
	ChangeSubmissionStatus(requester *models.Player, submissionID string, status models.SubmissionStatus) (*models.SubmissionView, error)
}

type ChangeSubmissionStatusProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeSubmissionStatusProcessor) New(gameRepository repositories.GameRepository) *ChangeSubmissionStatusProcessor {
	return &ChangeSubmissionStatusProcessor{
		gameRepository,
	}
}

// ChangeSubmissionStatus lets staff approve, reject or return a submission to the queue.
// When an admin rejects a submission it is vetoed: staff cannot bring it back, only an
// admin can. Played submissions belong to their round and cannot be changed.
func (P *ChangeSubmissionStatusProcessor) ChangeSubmissionStatus(requester *models.Player, submissionID string, status models.SubmissionStatus) (*models.SubmissionView, error) {
	if status != models.SUBMISSION_STATUS_PENDING &&
		status != models.SUBMISSION_STATUS_APPROVED &&
		status != models.SUBMISSION_STATUS_REJECTED {
		return nil, fmt.Errorf("%w: unknown status", ErrInvalid)
	}

	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(requester.UUID)
	if player == nil || (player.Class != models.PLAYER_CLASS_STAFF && player.Class != models.PLAYER_CLASS_ADMIN) {
		return nil, ErrForbidden
	}

	submission := game.GetSubmission(submissionID)
	if submission == nil {
		return nil, fmt.Errorf("%w: submission not found", ErrInvalid)
	}

	if submission.Status == models.SUBMISSION_STATUS_PLAYED {
		return nil, fmt.Errorf("%w: submission has already been played", ErrInvalid)
	}

	isAdmin := player.Class == models.PLAYER_CLASS_ADMIN
	if submission.Vetoed && !isAdmin {
		return nil, fmt.Errorf("%w: vetoed by the admin", ErrForbidden)
	}

	submission.Status = status
	if isAdmin {
		submission.Vetoed = status == models.SUBMISSION_STATUS_REJECTED
	}
	if status == models.SUBMISSION_STATUS_REJECTED && game.PinnedTopicID == submission.ID {
		game.PinnedTopicID = ""
	}
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game.SubmissionViewFor(submission, player), err
	}

	err = messages.SubmissionUpdatedMessage(game, submission)
	if err != nil {
		return game.SubmissionViewFor(submission, player), err
	}

	return game.SubmissionViewFor(submission, player), nil
}
