package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/requests"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
)

type ToggleSubmitterNamesController struct {
	controller.HTTP
	processor processors.ToggleSubmitterNamesProcessorInterface
}

func (C ToggleSubmitterNamesController) New(processor processors.ToggleSubmitterNamesProcessorInterface) *ToggleSubmitterNamesController {
	c := new(ToggleSubmitterNamesController)
	c.processor = processor
	return c
}

func (C *ToggleSubmitterNamesController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.ToggleSubmitterNames(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *ToggleSubmitterNamesController) ToggleSubmitterNames(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	toggleRequest := new(requests.ToggleRequest)
	err = fiberhelper.BindJSON(c, toggleRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ToggleSubmitterNames(player, toggleRequest.Enabled)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
