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

type ChangeCharacterRaceController struct {
	controller.HTTP
	processor processors.ChangeCharacterRaceProcessorInterface
}

func (C ChangeCharacterRaceController) New(processor processors.ChangeCharacterRaceProcessorInterface) *ChangeCharacterRaceController {
	c := new(ChangeCharacterRaceController)
	c.processor = processor

	return c
}

func (C *ChangeCharacterRaceController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeCharacterRace(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeCharacterRaceController) ChangeCharacterRace(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeCharacterRaceRequest := &struct {
		CharacterRace string
	}{
		"",
	}
	err = c.BindJSON(changeCharacterRaceRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)

	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeCharacterRace(player.RoomKey, uint(index), changeCharacterRaceRequest.CharacterRace)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
