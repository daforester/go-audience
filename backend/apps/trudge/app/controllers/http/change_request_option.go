package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChangeRequestOptionController struct {
	controller.HTTP
	processor     processors.ChangeRequestOptionProcessorInterface
	codeProcessor processors.SetGameCodeProcessorInterface
}

func (C ChangeRequestOptionController) New(processor processors.ChangeRequestOptionProcessorInterface) *ChangeRequestOptionController {
	c := new(ChangeRequestOptionController)
	c.processor = processor

	return c
}

func (C *ChangeRequestOptionController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeRequestOption(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeRequestOptionController) ChangeRequestOption(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeOptionRequest := &struct {
		Option string
	}{
		"",
	}
	err = c.BindJSON(changeOptionRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeRequestOption(player.RoomKey, uint(index), changeOptionRequest.Option)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
