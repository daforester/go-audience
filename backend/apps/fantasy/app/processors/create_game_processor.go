package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/libs/strgen"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CreateGameProcessorInterface interface {
	CreateGame(adminName string) (*models.GameView, *models.Player, error)
}

type CreateGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P CreateGameProcessor) New(gameRepository repositories.GameRepository) *CreateGameProcessor {
	return &CreateGameProcessor{
		gameRepository,
	}
}

func (P *CreateGameProcessor) CreateGame(adminName string) (*models.GameView, *models.Player, error) {
	adminName = CleanName(adminName)
	if adminName == "" {
		return nil, nil, fmt.Errorf("%w: name required", ErrInvalid)
	}

	admin := new(models.Player)
	id := uuid.New()
	admin.UUID = id.String()
	admin.Class = models.PLAYER_CLASS_ADMIN
	admin.Name = adminName

	game := models.NewGame()

	_ = P.gameRepository.Lock(nil)
	defer P.gameRepository.Unlock(nil)
	eGame := new(models.Game)
	for eGame != nil {
		game.Key = P.GenerateRoomKey()
		eGame = P.gameRepository.GetByKey(game.Key)
	}
	admin.RoomKey = game.Key
	game.AddPlayer(admin)
	game.Created = time.Now()
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, nil, err
	}

	return game.ViewFor(admin), admin.Copy(), nil
}

func (P *CreateGameProcessor) GenerateRoomKey() string {
	s, _ := strgen.RandString(5)
	return s
}
