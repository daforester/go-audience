package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

type GameRepository interface {
	Delete(game *models.Game) error
	GetByKey(key string) *models.Game
	Save(game *models.Game) error
	Lock(game *models.Game) error
	Unlock(game *models.Game) error
}
