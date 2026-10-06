package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChangeGamePasswordController struct {
	controller.HTTP
	processor processors.ChangeGamePasswordProcessorInterface
}

func (C ChangeGamePasswordController) New(processor processors.ChangeGamePasswordProcessorInterface) *ChangeGamePasswordController {
	c := new(ChangeGamePasswordController)
	c.processor = processor

	return c
}

func (C *ChangeGamePasswordController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeGamePassword(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeGamePasswordController) ChangeGamePassword(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	changeGamePasswordRequest := &struct {
		Password string
	}{
		"",
	}
	err = c.BindJSON(changeGamePasswordRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeGamePassword(changeGamePasswordRequest.Password)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
