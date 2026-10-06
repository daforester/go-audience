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

type StartRoundController struct {
	controller.HTTP
	processor processors.StartRoundProcessorInterface
}

func (C StartRoundController) New(processor processors.StartRoundProcessorInterface) *StartRoundController {
	c := new(StartRoundController)
	c.processor = processor
	return c
}

func (C *StartRoundController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.StartRound(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *StartRoundController) StartRound(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	startRoundRequest := new(requests.StartRoundRequest)
	err = fiberhelper.BindJSON(c, startRoundRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.StartRound(player, startRoundRequest.SubmissionID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
