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

type RemoveCharacterEquipmentController struct {
	controller.HTTP
	processor processors.RemoveCharacterEquipmentProcessorInterface
}

func (C RemoveCharacterEquipmentController) New(processor processors.RemoveCharacterEquipmentProcessorInterface) *RemoveCharacterEquipmentController {
	c := new(RemoveCharacterEquipmentController)
	c.processor = processor

	return c
}

func (C *RemoveCharacterEquipmentController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.RemoveCharacterEquipment(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *RemoveCharacterEquipmentController) RemoveCharacterEquipment(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	playerIndexStr := c.Param("playerIndex")
	playerIndex, err := strconv.ParseUint(playerIndexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	equipmentIndexStr := c.Param("equipmentIndex")
	equipmentIndex, err := strconv.ParseUint(equipmentIndexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.RemoveCharacterEquipment(player.RoomKey, uint(playerIndex), uint(equipmentIndex))

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
