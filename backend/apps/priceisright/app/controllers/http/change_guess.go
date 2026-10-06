package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChangeGuessController struct {
	controller.HTTP
	processor processors.ChangeGuessProcessorInterface
}

func (C ChangeGuessController) New(processor processors.ChangeGuessProcessorInterface) *ChangeGuessController {
	c := new(ChangeGuessController)
	c.processor = processor

	return c
}

func (C *ChangeGuessController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeGuess(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeGuessController) ChangeGuess(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeGuess := &struct {
		Round int64
		Value int64
	}{
		0,
		0,
	}
	err = c.BindJSON(changeGuess)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeGuess(uint(changeGuess.Round), uint(index), changeGuess.Value)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
