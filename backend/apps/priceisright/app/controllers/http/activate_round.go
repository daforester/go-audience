package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActivateRoundController struct {
	controller.HTTP
	processor processors.ActivateRoundProcessor
}

func (C ActivateRoundController) New(processor processors.ActivateRoundProcessor) *ActivateRoundController {
	c := new(ActivateRoundController)
	c.processor = processor
	return c
}

func (C *ActivateRoundController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ActivateRound(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ActivateRoundController) ActivateRound(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader to kick players")
	}

	ActivateRoundRequest := &struct {
		RoundNumber int
	}{
		-1,
	}
	err = c.BindJSON(ActivateRoundRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	round, err := C.processor.ActivateRound(ActivateRoundRequest.RoundNumber)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, "game not found")
	}

	return controller.StandardResponse(struct {
		Round *models.Round
	}{
		round,
	})
}
