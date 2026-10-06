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

type ReactTopicController struct {
	controller.HTTP
	processor processors.ReactTopicProcessorInterface
}

func (C ReactTopicController) New(processor processors.ReactTopicProcessorInterface) *ReactTopicController {
	c := new(ReactTopicController)
	c.processor = processor
	return c
}

func (C *ReactTopicController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.ReactTopic(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *ReactTopicController) ReactTopic(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	reactTopicRequest := new(requests.ReactTopicRequest)
	err = fiberhelper.BindJSON(c, reactTopicRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	round, err := C.processor.ReactTopic(player, reactTopicRequest.Reaction)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Round *models.RoundView
	}{
		round,
	})
}
