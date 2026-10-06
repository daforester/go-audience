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

type SetGameCodeController struct {
	controller.HTTP
	processor processors.SetGameCodeProcessorInterface
}

func (C SetGameCodeController) New(processor processors.SetGameCodeProcessorInterface) *SetGameCodeController {
	c := new(SetGameCodeController)
	c.processor = processor
	return c
}

func (C *SetGameCodeController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.SetGameCode(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *SetGameCodeController) SetGameCode(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	setGameCodeRequest := new(requests.SetGameCodeRequest)
	err = fiberhelper.BindJSON(c, setGameCodeRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.SetGameCode(player, setGameCodeRequest.Code)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
