package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/libs/strgen"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type TeamProcessorInterface interface {
	AddTeam(requester *models.Player, name string) (*models.GameView, error)
	ChangeTeam(requester *models.Player, teamID string, name string) (*models.GameView, error)
	RemoveTeam(requester *models.Player, teamID string) (*models.GameView, error)
}

type TeamProcessor struct {
	gameRepository repositories.GameRepository
}

func (P TeamProcessor) New(gameRepository repositories.GameRepository) *TeamProcessor {
	return &TeamProcessor{
		gameRepository,
	}
}

func (P *TeamProcessor) AddTeam(requester *models.Player, name string) (*models.GameView, error) {
	return P.changeTeams(requester, func(game *models.Game) error {
		if game.RoundInProgress() {
			return fmt.Errorf("%w: teams cannot change during a round", ErrWrongPhase)
		}
		if len(game.Teams) >= viper.GetInt("MAX_TEAMS") {
			return fmt.Errorf("%w: at most %d teams", ErrInvalid, viper.GetInt("MAX_TEAMS"))
		}

		name, err := P.teamName(game, "", name)
		if err != nil {
			return err
		}

		id, _ := strgen.RandString(8)
		for game.GetTeam(id) != nil {
			id, _ = strgen.RandString(8)
		}

		game.AddTeam(&models.Team{ID: id, Name: name})
		return nil
	})
}

func (P *TeamProcessor) ChangeTeam(requester *models.Player, teamID string, name string) (*models.GameView, error) {
	return P.changeTeams(requester, func(game *models.Game) error {
		team := game.GetTeam(teamID)
		if team == nil {
			return fmt.Errorf("%w: team not found", ErrInvalid)
		}

		name, err := P.teamName(game, teamID, name)
		if err != nil {
			return err
		}

		team.Name = name
		return nil
	})
}

func (P *TeamProcessor) RemoveTeam(requester *models.Player, teamID string) (*models.GameView, error) {
	return P.changeTeams(requester, func(game *models.Game) error {
		if game.RoundInProgress() {
			return fmt.Errorf("%w: teams cannot change during a round", ErrWrongPhase)
		}
		if game.GetTeam(teamID) == nil {
			return fmt.Errorf("%w: team not found", ErrInvalid)
		}

		directors := game.GetTeamDirectors(teamID)
		game.RemoveTeam(teamID)

		// Directors of the removed team are back in the audience
		wsCollection := sockets.GetCollection(game.Key)
		for _, p := range directors {
			for _, sock := range wsCollection.HasUUID(p.UUID) {
				_, _ = sock.SwapTags(models.FantasyPlayerClassesStrings, []string{string(p.Class)})
			}
			_ = messages.PlayerUpdateMessage(p)
		}
		return nil
	})
}

func (P *TeamProcessor) teamName(game *models.Game, teamID string, name string) (string, error) {
	name = CleanName(name)
	if name == "" {
		return "", fmt.Errorf("%w: team name required", ErrInvalid)
	}
	for _, t := range game.Teams {
		if t.ID != teamID && strings.EqualFold(t.Name, name) {
			return "", fmt.Errorf("%w: another team already has that name", ErrInvalid)
		}
	}
	return name, nil
}

func (P *TeamProcessor) changeTeams(requester *models.Player, change func(game *models.Game) error) (*models.GameView, error) {
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

	err := change(game)
	if err != nil {
		return nil, err
	}

	game.Updated = time.Now()
	err = P.gameRepository.Save(game)

	if err != nil {
		return messages.View(game, admin), err
	}

	err = messages.GameUpdateMessage(game)
	if err != nil {
		return messages.View(game, admin), err
	}

	return messages.View(game, admin), nil
}
