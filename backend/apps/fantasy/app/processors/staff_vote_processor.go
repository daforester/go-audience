package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type StaffVoteProcessorInterface interface {
	StaffVote(requester *models.Player, submissionID string) (*models.SubmissionView, error)
}

type StaffVoteProcessor struct {
	gameRepository repositories.GameRepository
}

func (P StaffVoteProcessor) New(gameRepository repositories.GameRepository) *StaffVoteProcessor {
	return &StaffVoteProcessor{
		gameRepository,
	}
}

// StaffVote toggles the requester's upvote on a submission still in the queue
func (P *StaffVoteProcessor) StaffVote(requester *models.Player, submissionID string) (*models.SubmissionView, error) {
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

	if submission.Status != models.SUBMISSION_STATUS_PENDING && submission.Status != models.SUBMISSION_STATUS_APPROVED {
		return nil, fmt.Errorf("%w: submission is no longer in the queue", ErrInvalid)
	}

	submission.ToggleStaffVote(player.UUID)
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
