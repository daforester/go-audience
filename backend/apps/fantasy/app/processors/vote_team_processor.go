package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type VoteTeamProcessorInterface interface {
	VoteTeam(requester *models.Player, teamID string) (*models.RoundView, error)
}

type VoteTeamProcessor struct {
	gameRepository repositories.GameRepository
}

func (P VoteTeamProcessor) New(gameRepository repositories.GameRepository) *VoteTeamProcessor {
	return &VoteTeamProcessor{
		gameRepository,
	}
}

// VoteTeam records an audience member's vote for the team with the best answer. Each
// member has one vote per round, which they may change until voting closes.
func (P *VoteTeamProcessor) VoteTeam(requester *models.Player, teamID string) (*models.RoundView, error) {
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

	if game.GameCode != models.GAME_CODE_AUDIENCE_VOTE {
		return nil, fmt.Errorf("%w: voting is not open", ErrWrongPhase)
	}

	round := game.GetCurrentRound()
	if round == nil || !round.HasTeam(teamID) {
		return nil, fmt.Errorf("%w: team is not in this round", ErrInvalid)
	}

	round.Vote(player.UUID, teamID)
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, err
	}

	messages.VoteTallyMessage(game, round)

	return game.RoundViewFor(player), nil
}
