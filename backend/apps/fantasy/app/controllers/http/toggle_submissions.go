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

type ToggleSubmissionsController struct {
	controller.HTTP
	processor processors.ToggleSubmissionsProcessorInterface
}

func (C ToggleSubmissionsController) New(processor processors.ToggleSubmissionsProcessorInterface) *ToggleSubmissionsController {
	c := new(ToggleSubmissionsController)
	c.processor = processor
	return c
}

func (C *ToggleSubmissionsController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.ToggleSubmissions(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *ToggleSubmissionsController) ToggleSubmissions(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	toggleRequest := new(requests.ToggleRequest)
	err = fiberhelper.BindJSON(c, toggleRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ToggleSubmissions(player, toggleRequest.Enabled)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
