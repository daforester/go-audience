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

type AddCharacterEquipmentController struct {
	controller.HTTP
	processor processors.AddCharacterEquipmentProcessorInterface
}

func (C AddCharacterEquipmentController) New(processor processors.AddCharacterEquipmentProcessorInterface) *AddCharacterEquipmentController {
	c := new(AddCharacterEquipmentController)
	c.processor = processor

	return c
}

func (C *AddCharacterEquipmentController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AddCharacterEquipment(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AddCharacterEquipmentController) AddCharacterEquipment(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	addCharacterEquipmentRequest := &struct {
		Equipment string
	}{
		"",
	}
	err = c.BindJSON(addCharacterEquipmentRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	playerIndexStr := c.Param("playerIndex")
	playerIndex, err := strconv.ParseUint(playerIndexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.AddCharacterEquipment(player.RoomKey, uint(playerIndex), addCharacterEquipmentRequest.Equipment)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
