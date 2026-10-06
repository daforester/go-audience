package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EndGameController struct {
	controller.HTTP
	processor processors.QuitGameProcessor
}

func (C EndGameController) New(processor processors.QuitGameProcessor) *EndGameController {
	c := new(EndGameController)
	c.processor = processor
	return c
}

func (C *EndGameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.EndGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *EndGameController) EndGame(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	err = C.processor.QuitGame(player)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		GameCode int
		Game     *models.Game
		Player   *models.Player
	}{
		0,
		nil,
		nil,
	})
}
