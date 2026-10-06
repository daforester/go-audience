package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChangePlayerExperienceController struct {
	controller.HTTP
	processor processors.ChangePlayerExperienceProcessorInterface
}

func (C ChangePlayerExperienceController) New(processor processors.ChangePlayerExperienceProcessorInterface) *ChangePlayerExperienceController {
	c := new(ChangePlayerExperienceController)
	c.processor = processor

	return c
}

func (C *ChangePlayerExperienceController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangePlayerExperience(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangePlayerExperienceController) ChangePlayerExperience(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	ChangePlayerExperienceRequest := &struct {
		PlayerUUID string
		Adjustment int
	}{
		"",
		0,
	}
	err = c.BindJSON(ChangePlayerExperienceRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangePlayerExperience(player.RoomKey, ChangePlayerExperienceRequest.PlayerUUID, ChangePlayerExperienceRequest.Adjustment)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
