package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type ReactTopicProcessorInterface interface {
	ReactTopic(requester *models.Player, reaction int) (*models.RoundView, error)
}

type ReactTopicProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ReactTopicProcessor) New(gameRepository repositories.GameRepository) *ReactTopicProcessor {
	return &ReactTopicProcessor{
		gameRepository,
	}
}

// ReactTopic records an audience member's +1 / -1 on the live topic while teams present.
// 0 withdraws the reaction.
func (P *ReactTopicProcessor) ReactTopic(requester *models.Player, reaction int) (*models.RoundView, error) {
	if reaction < -1 || reaction > 1 {
		return nil, fmt.Errorf("%w: reaction must be -1, 0 or 1", ErrInvalid)
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

	if game.GameCode != models.GAME_CODE_RESPONSES {
		return nil, fmt.Errorf("%w: reactions are open while teams present", ErrWrongPhase)
	}

	round := game.GetCurrentRound()
	if round == nil {
		return nil, fmt.Errorf("%w: no topic", ErrWrongPhase)
	}

	topic := game.GetSubmission(round.SubmissionID)
	if topic == nil {
		return nil, fmt.Errorf("%w: no topic", ErrWrongPhase)
	}

	topic.SetReaction(player.UUID, reaction)
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, err
	}

	messages.ReactionTallyMessage(game, topic)

	return game.RoundViewFor(player), nil
}
