package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChangePlayerHealthController struct {
	controller.HTTP
	processor processors.ChangePlayerHealthProcessorInterface
}

func (C ChangePlayerHealthController) New(processor processors.ChangePlayerHealthProcessorInterface) *ChangePlayerHealthController {
	c := new(ChangePlayerHealthController)
	c.processor = processor

	return c
}

func (C *ChangePlayerHealthController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangePlayerHealth(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangePlayerHealthController) ChangePlayerHealth(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	ChangePlayerHealthRequest := &struct {
		PlayerUUID string
		Adjustment int
	}{
		"",
		0,
	}
	err = c.BindJSON(ChangePlayerHealthRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangePlayerHealth(player.RoomKey, ChangePlayerHealthRequest.PlayerUUID, ChangePlayerHealthRequest.Adjustment)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
