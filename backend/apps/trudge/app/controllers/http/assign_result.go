package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/requests"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssignResultController struct {
	controller.HTTP
	processor processors.AssignResultProcessorInterface
}

func (C AssignResultController) New(processor processors.AssignResultProcessorInterface) *AssignResultController {
	c := new(AssignResultController)
	c.processor = processor
	return c
}

func (C *AssignResultController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AssignResult(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AssignResultController) AssignResult(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader")
	}

	assignResultRequest := new(requests.AssignResultRequest)
	err = c.BindJSON(assignResultRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.AssignResult(player.RoomKey, assignResultRequest.Property)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
