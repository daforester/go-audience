package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuitGameController struct {
	controller.HTTP
	processor processors.QuitGameProcessor
}

func (C QuitGameController) New(processor processors.QuitGameProcessor) *QuitGameController {
	c := new(QuitGameController)
	c.processor = processor
	return c
}

func (C *QuitGameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.QuitGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *QuitGameController) QuitGame(c *gin.Context) *controller.Response {
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
