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

type ChangePlayerBadgeController struct {
	controller.HTTP
	processor processors.ChangePlayerBadgeProcessorInterface
}

func (C ChangePlayerBadgeController) New(processor processors.ChangePlayerBadgeProcessorInterface) *ChangePlayerBadgeController {
	c := new(ChangePlayerBadgeController)
	c.processor = processor

	return c
}

func (C *ChangePlayerBadgeController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangePlayerBadge(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangePlayerBadgeController) ChangePlayerBadge(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changePlayerBadgeRequest := &struct {
		BadgeNumber string
	}{
		"",
	}
	err = c.BindJSON(changePlayerBadgeRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangePlayerBadge(uint(index), changePlayerBadgeRequest.BadgeNumber)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
