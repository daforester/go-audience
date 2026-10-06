package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type SetGameCodeProcessorInterface interface {
	SetGameCode(requester *models.Player, code int) (*models.GameView, error)
}

type SetGameCodeProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SetGameCodeProcessor) New(gameRepository repositories.GameRepository) *SetGameCodeProcessor {
	return &SetGameCodeProcessor{
		gameRepository,
	}
}

func (P *SetGameCodeProcessor) SetGameCode(requester *models.Player, code int) (*models.GameView, error) {
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

	if code == models.GAME_CODE_TOPIC_REVEAL {
		return nil, fmt.Errorf("%w: start a round to reveal a topic", ErrWrongPhase)
	}

	if !models.CanTransition(game.GameCode, code) {
		return nil, fmt.Errorf("%w: cannot go from %d to %d", ErrWrongPhase, game.GameCode, code)
	}

	round := game.GetCurrentRound()

	switch code {
	case models.GAME_CODE_LOBBY:
		game.Reset()
	case models.GAME_CODE_INTAKE:
		if game.GameCode == models.GAME_CODE_LOBBY {
			game.SubmissionsEnabled = true
		}
	case models.GAME_CODE_RESPONSES:
		if round != nil && len(round.TeamOrder) > 0 {
			round.PresentingTeam = round.TeamOrder[0]
		}
	case models.GAME_CODE_AUDIENCE_VOTE:
		if round != nil {
			round.PresentingTeam = ""
		}
	case models.GAME_CODE_ROUND_RESULT:
		game.DecideRound()
	}

	game.GameCode = code
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	messages.SetGameCodeMessage(game.Key, code)

	err = messages.GameUpdateMessage(game)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
