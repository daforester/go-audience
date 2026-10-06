package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type PinTopicProcessorInterface interface {
	PinTopic(requester *models.Player, submissionID string) (*models.GameView, error)
}

type PinTopicProcessor struct {
	gameRepository repositories.GameRepository
}

func (P PinTopicProcessor) New(gameRepository repositories.GameRepository) *PinTopicProcessor {
	return &PinTopicProcessor{
		gameRepository,
	}
}

// PinTopic lets the admin override the staff's top pick: the pinned submission is the
// next topic revealed. An empty ID removes the pin.
func (P *PinTopicProcessor) PinTopic(requester *models.Player, submissionID string) (*models.GameView, error) {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	admin := game.GetPlayerByUUID(requester.UUID)
	if admin == nil || admin.Class != models.PLAYER_CLASS_ADMIN {
		return nil, ErrForbidden
	}

	if submissionID != "" {
		submission := game.GetSubmission(submissionID)
		if submission == nil {
			return nil, fmt.Errorf("%w: submission not found", ErrInvalid)
		}
		if !submission.Eligible() {
			return nil, fmt.Errorf("%w: submission has been %s", ErrInvalid, submission.Status)
		}
	}

	game.PinnedTopicID = submissionID
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	err = messages.QueueUpdatedMessage(game)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
