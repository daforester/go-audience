package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"github.com/gofiber/fiber/v3"
)

type RejoinGameController struct {
	controller.HTTP
	processor processors.GetGameProcessorInterface
}

func (C RejoinGameController) New(processor processors.GetGameProcessorInterface) *RejoinGameController {
	c := new(RejoinGameController)
	c.processor = processor
	return c
}

func (C *RejoinGameController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.RejoinGame(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *RejoinGameController) RejoinGame(c fiber.Ctx) *controller.Response {
	requester, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	game, player, err := C.processor.GetGame(requester)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	return controller.StandardResponse(struct {
		Game   *models.GameView
		Player *models.Player
	}{
		game,
		player,
	})
}
