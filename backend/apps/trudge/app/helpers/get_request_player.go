package helpers

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/middleware"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetRequestPlayer(c *gin.Context) (*models.Player, error) {
	playerInterface, exists := c.Get("Player")
	if !exists {
		err := middleware.ProcessPlayerToken(c)
		if err != nil {
			return nil, errors.New("player token data not provided")
		}
		playerInterface, exists = c.Get("Player")
		if !exists {
			return nil, errors.New("player token data not parsed")
		}
	}

	player := playerInterface.(*models.Player)

	return player, nil
}
