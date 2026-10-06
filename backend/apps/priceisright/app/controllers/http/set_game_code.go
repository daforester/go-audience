package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SetGameCodeController struct {
	controller.HTTP
	processor processors.SetGameCodeProcessor
}

func (C SetGameCodeController) New(processor processors.SetGameCodeProcessor) *SetGameCodeController {
	c := new(SetGameCodeController)
	c.processor = processor
	return c
}

func (C *SetGameCodeController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.SetGameCode(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *SetGameCodeController) SetGameCode(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader to kick players")
	}

	setGameCodeRequest := &struct {
		Code int
	}{
		0,
	}
	err = c.BindJSON(setGameCodeRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.SetGameCode(setGameCodeRequest.Code)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, "game not found")
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
