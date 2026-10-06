package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
)

type GameRepository interface {
	Get() *models.Game
	Save(*models.Game)
	Lock()
	Unlock()
}
