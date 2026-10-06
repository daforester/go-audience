package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"github.com/gofiber/fiber/v3"
)

type GetGameController struct {
	controller.HTTP
	processor processors.GetGameProcessorInterface
}

func (C GetGameController) New(processor processors.GetGameProcessorInterface) *GetGameController {
	c := new(GetGameController)
	c.processor = processor
	return c
}

func (C *GetGameController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.GetGame(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *GetGameController) GetGame(c fiber.Ctx) *controller.Response {
	requester, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	game, _, err := C.processor.GetGame(requester)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
