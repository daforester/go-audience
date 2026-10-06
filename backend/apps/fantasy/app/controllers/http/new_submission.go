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

type NewSubmissionController struct {
	controller.HTTP
	processor processors.NewSubmissionProcessorInterface
}

func (C NewSubmissionController) New(processor processors.NewSubmissionProcessorInterface) *NewSubmissionController {
	c := new(NewSubmissionController)
	c.processor = processor
	return c
}

func (C *NewSubmissionController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.NewSubmission(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *NewSubmissionController) NewSubmission(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	newSubmissionRequest := new(requests.NewSubmissionRequest)
	err = fiberhelper.BindJSON(c, newSubmissionRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	submission, err := C.processor.NewSubmission(player, newSubmissionRequest.Kind, newSubmissionRequest.Text)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Submission *models.SubmissionView
	}{
		submission,
	})
}
