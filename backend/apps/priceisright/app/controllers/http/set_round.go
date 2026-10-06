package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SetRoundController struct {
	controller.HTTP
	processor processors.SetRoundProcessor
}

func (C SetRoundController) New(processor processors.SetRoundProcessor) *SetRoundController {
	c := new(SetRoundController)
	c.processor = processor
	return c
}

func (C *SetRoundController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.SetRound(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *SetRoundController) SetRound(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader to kick players")
	}

	setRoundRequest := &struct {
		RoundNumber int
	}{
		0,
	}
	err = c.BindJSON(setRoundRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.SetRound(setRoundRequest.RoundNumber)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, "game not found")
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
