package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssignRaceController struct {
	controller.HTTP
	processor processors.AssignRaceProcessorInterface
}

func (C AssignRaceController) New(processor processors.AssignRaceProcessorInterface) *AssignRaceController {
	c := new(AssignRaceController)
	c.processor = processor
	return c
}

func (C *AssignRaceController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AssignRace(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AssignRaceController) AssignRace(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_BARD &&
		player.Class != models.PLAYER_CLASS_WARRIOR &&
		player.Class != models.PLAYER_CLASS_WIZARD &&
		player.Class != models.PLAYER_CLASS_GOD {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be party member or god")
	}

	assignRaceRequest := &struct {
		CharacterRace string
	}{
		"",
	}
	err = c.BindJSON(assignRaceRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	p, err := C.processor.AssignRace(player.RoomKey, player.UUID, assignRaceRequest.CharacterRace)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Player *models.Player
	}{
		p,
	})
}
