package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"github.com/gofiber/fiber/v3"
)

type QuitGameController struct {
	controller.HTTP
	processor processors.QuitGameProcessorInterface
}

func (C QuitGameController) New(processor processors.QuitGameProcessorInterface) *QuitGameController {
	c := new(QuitGameController)
	c.processor = processor
	return c
}

func (C *QuitGameController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.QuitGame(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *QuitGameController) QuitGame(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	err = C.processor.QuitGame(player)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		GameCode int
		Game     *models.GameView
		Player   *models.Player
	}{
		0,
		nil,
		nil,
	})
}
