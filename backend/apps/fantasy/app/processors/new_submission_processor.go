package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

// Upper bounds so one room, or one player, cannot fill the server's memory
const MaxSubmissionsPerGame = 2000
const MaxSubmissionsPerPlayer = 50

type NewSubmissionProcessorInterface interface {
	NewSubmission(requester *models.Player, kind models.SubmissionKind, text string) (*models.SubmissionView, error)
}

type NewSubmissionProcessor struct {
	gameRepository repositories.GameRepository
}

func (P NewSubmissionProcessor) New(gameRepository repositories.GameRepository) *NewSubmissionProcessor {
	return &NewSubmissionProcessor{
		gameRepository,
	}
}

func (P *NewSubmissionProcessor) NewSubmission(requester *models.Player, kind models.SubmissionKind, text string) (*models.SubmissionView, error) {
	if !models.ValidSubmissionKind(kind) {
		return nil, fmt.Errorf("%w: unknown kind of submission", ErrInvalid)
	}

	text = CleanText(text, viper.GetInt("MAX_SUBMISSION_LENGTH"))
	if text == "" {
		return nil, fmt.Errorf("%w: submission is empty", ErrInvalid)
	}

	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return nil, ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(requester.UUID)
	if player == nil || player.Class != models.PLAYER_CLASS_AUDIENCE {
		return nil, ErrForbidden
	}

	if !game.SubmissionsEnabled {
		return nil, fmt.Errorf("%w: submissions are closed", ErrWrongPhase)
	}

	if game.LastSubmission == nil {
		game.LastSubmission = make(map[string]time.Time)
	}
	cooldown := time.Duration(viper.GetInt("SUBMISSION_COOLDOWN_SECONDS")) * time.Second
	if last, ok := game.LastSubmission[player.UUID]; ok && time.Since(last) < cooldown {
		return nil, ErrTooSoon
	}

	if len(game.Submissions) >= MaxSubmissionsPerGame {
		return nil, fmt.Errorf("%w: the submission box is full", ErrInvalid)
	}
	mine := 0
	for _, s := range game.Submissions {
		if s.AuthorUUID == player.UUID {
			mine++
		}
	}
	if mine >= MaxSubmissionsPerPlayer {
		return nil, fmt.Errorf("%w: you have sent the most submissions allowed", ErrInvalid)
	}

	submission := models.NewSubmission(uuid.New().String(), player, kind, text)
	game.AddSubmission(submission)
	game.LastSubmission[player.UUID] = time.Now()
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, err
	}

	err = messages.NewSubmissionMessage(game, submission)
	if err != nil {
		return game.SubmissionViewFor(submission, player), err
	}

	return game.SubmissionViewFor(submission, player), nil
}
