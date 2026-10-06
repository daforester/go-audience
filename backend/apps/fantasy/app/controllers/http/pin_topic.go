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

type PinTopicController struct {
	controller.HTTP
	processor processors.PinTopicProcessorInterface
}

func (C PinTopicController) New(processor processors.PinTopicProcessorInterface) *PinTopicController {
	c := new(PinTopicController)
	c.processor = processor
	return c
}

func (C *PinTopicController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.PinTopic(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *PinTopicController) PinTopic(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	pinTopicRequest := new(requests.PinTopicRequest)
	err = fiberhelper.BindJSON(c, pinTopicRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.PinTopic(player, pinTopicRequest.SubmissionID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
