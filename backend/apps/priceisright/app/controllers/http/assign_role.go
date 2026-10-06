package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/requests"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssignRoleController struct {
	controller.HTTP
	processor processors.AssignRoleProcessorInterface
}

func (C AssignRoleController) New(processor processors.AssignRoleProcessorInterface) *AssignRoleController {
	c := new(AssignRoleController)
	c.processor = processor
	return c
}

func (C *AssignRoleController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AssignRole(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AssignRoleController) AssignRole(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader")
	}

	assignPlayerRequest := new(requests.AssignPlayerRequest)
	err = c.BindJSON(assignPlayerRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	p, err := C.processor.AssignPlayer(assignPlayerRequest.UUID, assignPlayerRequest.PlayerClass)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Player *models.Player
	}{
		p,
	})
}
