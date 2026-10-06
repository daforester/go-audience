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

type ChangeCharacterEquipmentController struct {
	controller.HTTP
	processor processors.ChangeCharacterEquipmentProcessorInterface
}

func (C ChangeCharacterEquipmentController) New(processor processors.ChangeCharacterEquipmentProcessorInterface) *ChangeCharacterEquipmentController {
	c := new(ChangeCharacterEquipmentController)
	c.processor = processor

	return c
}

func (C *ChangeCharacterEquipmentController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeCharacterEquipment(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeCharacterEquipmentController) ChangeCharacterEquipment(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeCharacterEquipmentRequest := &struct {
		Equipment string
	}{
		"",
	}
	err = c.BindJSON(changeCharacterEquipmentRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
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

	game, err := C.processor.ChangeCharacterEquipment(player.RoomKey, uint(playerIndex), uint(equipmentIndex), changeCharacterEquipmentRequest.Equipment)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
