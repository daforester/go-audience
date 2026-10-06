package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type SetPresentingTeamProcessorInterface interface {
	SetPresentingTeam(requester *models.Player, teamID string) error
}

type SetPresentingTeamProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SetPresentingTeamProcessor) New(gameRepository repositories.GameRepository) *SetPresentingTeamProcessor {
	return &SetPresentingTeamProcessor{
		gameRepository,
	}
}

func (P *SetPresentingTeamProcessor) SetPresentingTeam(requester *models.Player, teamID string) error {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	admin := game.GetPlayerByUUID(requester.UUID)
	if admin == nil || admin.Class != models.PLAYER_CLASS_ADMIN {
		return ErrForbidden
	}

	if game.GameCode != models.GAME_CODE_RESPONSES {
		return fmt.Errorf("%w: teams present during responses", ErrWrongPhase)
	}

	round := game.GetCurrentRound()
	if round == nil || !round.HasTeam(teamID) {
		return fmt.Errorf("%w: team is not in this round", ErrInvalid)
	}

	round.PresentingTeam = teamID
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return err
	}

	return messages.PresentingTeamMessage(game.Key, teamID)
}
