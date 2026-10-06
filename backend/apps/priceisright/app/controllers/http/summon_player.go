package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SummonPlayerController struct {
	controller.HTTP
	processor processors.SummonPlayerProcessor
}

func (C SummonPlayerController) New(processor processors.SummonPlayerProcessor) *SummonPlayerController {
	c := new(SummonPlayerController)
	c.processor = processor
	return c
}

func (C *SummonPlayerController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.SummonPlayer(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *SummonPlayerController) SummonPlayer(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader to summon players")
	}

	SummonPlayerRequest := &struct {
		Player *models.Player
		Value  int
	}{
		nil,
		0,
	}
	err = c.BindJSON(SummonPlayerRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	err = C.processor.SummonPlayer(SummonPlayerRequest.Player, SummonPlayerRequest.Value)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, "game not found")
	}

	return controller.StandardResponse(struct {
		done bool
	}{
		true,
	})
}
