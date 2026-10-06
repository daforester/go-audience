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

type ChangeCharacterNameController struct {
	controller.HTTP
	processor processors.ChangeCharacterNameProcessorInterface
}

func (C ChangeCharacterNameController) New(processor processors.ChangeCharacterNameProcessorInterface) *ChangeCharacterNameController {
	c := new(ChangeCharacterNameController)
	c.processor = processor

	return c
}

func (C *ChangeCharacterNameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeCharacterName(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeCharacterNameController) ChangeCharacterName(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeCharacterNameRequest := &struct {
		CharacterName string
	}{
		"",
	}
	err = c.BindJSON(changeCharacterNameRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeCharacterName(player.RoomKey, uint(index), changeCharacterNameRequest.CharacterName)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
