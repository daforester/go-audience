package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-backend-repository/component/libs/strgen"
	"github.com/google/uuid"
	"time"
)

type CreateGameProcessorInterface interface {
	CreateGame(leaderName string) (*models.Game, *models.Player, error)
}

type CreateGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P CreateGameProcessor) New(gameRepository repositories.GameRepository) *CreateGameProcessor {
	return &CreateGameProcessor{
		gameRepository,
	}
}

func (P *CreateGameProcessor) CreateGame(leaderName string) (*models.Game, *models.Player, error) {
	leader := new(models.Player)
	id := uuid.New()
	leader.UUID = id.String()
	leader.Class = models.PLAYER_CLASS_LEADER
	leader.Name = leaderName

	game := new(models.Game)
	game.GameCode = 20
	game.Created = time.Now()
	game.Updated = time.Now()

	_ = P.gameRepository.Lock(nil)
	defer P.gameRepository.Unlock(nil)
	eGame := new(models.Game)
	for eGame != nil {
		game.Key = P.GenerateRoomKey()
		eGame = P.gameRepository.GetByKey(game.Key)
	}
	leader.RoomKey = game.Key
	game.AddPlayer(leader)
	err := P.gameRepository.Save(game)

	if err != nil {
		return nil, nil, err
	}

	return game, leader, nil
}

func (P *CreateGameProcessor) GenerateRoomKey() string {
	s, _ := strgen.RandString(5)
	return s
}
