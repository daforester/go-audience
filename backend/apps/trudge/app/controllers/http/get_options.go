package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/middleware"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetOptionsController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C GetOptionsController) New(gameRepository repositories.GameRepository) *GetOptionsController {
	c := new(GetOptionsController)
	c.gameRepository = gameRepository
	return c
}

func (C *GetOptionsController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.GetOptions(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *GetOptionsController) GetOptions(c *gin.Context) *controller.Response {
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

	game := C.gameRepository.GetByKey(player.RoomKey)

	if game == nil {
		return controller.ErrorResponse(http.StatusUnauthorized, "game not found")
	}

	return controller.StandardResponse(struct {
		Options []string
	}{
		game.Options,
	})
}
