package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"time"
)

type AssignRoleProcessorInterface interface {
	AssignPlayer(admin *models.Player, playerUUID string, class models.PlayerClass, teamID string) (*models.Player, error)
}

type AssignRoleProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AssignRoleProcessor) New(gameRepository repositories.GameRepository) *AssignRoleProcessor {
	return &AssignRoleProcessor{
		gameRepository,
	}
}

func (P *AssignRoleProcessor) AssignPlayer(requester *models.Player, playerUUID string, class models.PlayerClass, teamID string) (*models.Player, error) {
	if !models.ValidPlayerClass(class) {
		return nil, fmt.Errorf("%w: unknown role", ErrInvalid)
	}

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

	player := game.GetPlayerByUUID(playerUUID)
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	if player.Class == models.PLAYER_CLASS_ADMIN && class != models.PLAYER_CLASS_ADMIN && len(game.GetAdmins()) == 1 {
		return nil, fmt.Errorf("%w: the game needs at least one admin", ErrInvalid)
	}

	if class == models.PLAYER_CLASS_DIRECTOR {
		if game.GetTeam(teamID) == nil {
			return nil, fmt.Errorf("%w: directors must be assigned to a team", ErrInvalid)
		}
	} else {
		teamID = ""
	}

	player.Class = class
	player.TeamID = teamID

	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return player.Copy(), err
	}

	err = P.swapSocketTags(player)
	if err != nil {
		return player.Copy(), err
	}

	err = messages.PlayerUpdateMessage(player)
	if err != nil {
		return player.Copy(), err
	}

	err = messages.GameUpdateMessage(game)
	if err != nil {
		return player.Copy(), err
	}

	return player.Copy(), nil
}

func (P *AssignRoleProcessor) swapSocketTags(player *models.Player) error {
	wsCollection := sockets.GetCollection(player.RoomKey)
	socks := wsCollection.HasUUID(player.UUID)
	for _, sock := range socks {
		_, err := sock.SwapTags(models.FantasyPlayerClassesStrings, []string{string(player.Class)})
		if err != nil {
			return err
		}
	}
	return nil
}
