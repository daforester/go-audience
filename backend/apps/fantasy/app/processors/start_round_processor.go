package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type StartRoundProcessorInterface interface {
	StartRound(requester *models.Player, submissionID string) (*models.GameView, error)
}

type StartRoundProcessor struct {
	gameRepository repositories.GameRepository
}

func (P StartRoundProcessor) New(gameRepository repositories.GameRepository) *StartRoundProcessor {
	return &StartRoundProcessor{
		gameRepository,
	}
}

// StartRound reveals the next topic: the given submission if the admin picked one,
// otherwise the approved submission with the most staff votes
func (P *StartRoundProcessor) StartRound(requester *models.Player, submissionID string) (*models.GameView, error) {
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

	if !models.CanTransition(game.GameCode, models.GAME_CODE_TOPIC_REVEAL) {
		return nil, fmt.Errorf("%w: finish the current round first", ErrWrongPhase)
	}

	if err := game.ReadyForRound(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalid, err.Error())
	}

	var topic *models.Submission
	if submissionID != "" {
		topic = game.GetSubmission(submissionID)
		if topic == nil {
			return nil, fmt.Errorf("%w: submission not found", ErrInvalid)
		}
		if topic.Status == models.SUBMISSION_STATUS_PLAYED || topic.Status == models.SUBMISSION_STATUS_REJECTED {
			return nil, fmt.Errorf("%w: submission has been %s", ErrInvalid, topic.Status)
		}
	} else {
		topic = game.NextTopic()
		if topic == nil {
			return nil, fmt.Errorf("%w: no approved submissions to play", ErrInvalid)
		}
	}

	game.StartRound(topic)
	game.GameCode = models.GAME_CODE_TOPIC_REVEAL
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	messages.SetGameCodeMessage(game.Key, game.GameCode)

	err = messages.GameUpdateMessage(game)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
