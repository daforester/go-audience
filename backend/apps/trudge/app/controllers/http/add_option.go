package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddOptionController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	processor      processors.AddOptionProcessorInterface
}

func (C AddOptionController) New(gameRepository repositories.GameRepository, processor processors.AddOptionProcessorInterface) *AddOptionController {
	c := new(AddOptionController)
	c.gameRepository = gameRepository
	c.processor = processor
	return c
}

func (C *AddOptionController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.AddOption(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *AddOptionController) AddOption(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_GOD {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be a god")
	}

	addOptionRequest := &struct {
		Option string
	}{
		"",
	}
	err = c.BindJSON(addOptionRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, err = C.processor.AddOption(player.RoomKey, player.UUID, addOptionRequest.Option)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Option string
	}{
		addOptionRequest.Option,
	})
}
