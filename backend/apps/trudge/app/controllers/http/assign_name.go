package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssignNameController struct {
	controller.HTTP
	processor processors.AssignNameProcessorInterface
}

func (C AssignNameController) New(processor processors.AssignNameProcessorInterface) *AssignNameController {
	c := new(AssignNameController)
	c.processor = processor
	return c
}

func (C *AssignNameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AssignName(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AssignNameController) AssignName(c *gin.Context) *controller.Response {
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

	assignNameRequest := &struct {
		CharacterName string
	}{
		"",
	}
	err = c.BindJSON(assignNameRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	p, err := C.processor.AssignName(player.RoomKey, player.UUID, assignNameRequest.CharacterName)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Player *models.Player
	}{
		p,
	})
}
