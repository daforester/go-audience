package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetGameController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C GetGameController) New(gameRepository repositories.GameRepository) *GetGameController {
	c := new(GetGameController)
	c.gameRepository = gameRepository
	return c
}

func (C *GetGameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.GetGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *GetGameController) GetGame(c *gin.Context) *controller.Response {
	playerInterface, exists := c.Get("Player")
	if !exists {
		err := middleware.ProcessPlayerToken(c)
		if err != nil {
			return controller.ErrorResponse(http.StatusUnauthorized, "player token data not provided")
		}
		playerInterface, exists = c.Get("Player")
		if !exists {
			return controller.ErrorResponse(http.StatusUnauthorized, "player token data not parsed")
		}
	}

	player := playerInterface.(*models.Player)

	if player.Class != models.PLAYER_CLASS_LEADER && player.Class != models.PLAYER_CLASS_DISPLAY {
		return controller.ErrorResponse(http.StatusUnauthorized, "access denied")
	}

	game := C.gameRepository.Get()

	if game == nil {
		return controller.ErrorResponse(http.StatusUnauthorized, "game not found")
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
