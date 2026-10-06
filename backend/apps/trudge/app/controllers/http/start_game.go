package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StartGameController struct {
	controller.HTTP
	processor processors.SetGameCodeProcessor
}

func (C StartGameController) New(processor processors.SetGameCodeProcessor) *StartGameController {
	c := new(StartGameController)
	c.processor = processor
	return c
}

func (C *StartGameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.StartGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *StartGameController) StartGame(c *gin.Context) *controller.Response {
	playerInterface, exists := c.Get("Player")
	if !exists {
		return controller.ErrorResponse(http.StatusUnauthorized, "player token data not parsed")
	}

	player := playerInterface.(*models.Player)

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "only leaders may start game")
	}

	game, err := C.processor.SetGameCode(player.RoomKey, 30)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, "game not found")
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
