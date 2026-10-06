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

type ChangePlayerNameController struct {
	controller.HTTP
	processor processors.ChangePlayerNameProcessorInterface
}

func (C ChangePlayerNameController) New(processor processors.ChangePlayerNameProcessorInterface) *ChangePlayerNameController {
	c := new(ChangePlayerNameController)
	c.processor = processor

	return c
}

func (C *ChangePlayerNameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangePlayerName(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangePlayerNameController) ChangePlayerName(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changePlayerNameRequest := &struct {
		Name string
	}{
		"",
	}
	err = c.BindJSON(changePlayerNameRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangePlayerName(uint(index), changePlayerNameRequest.Name)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
