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

type ChangeSubmissionStatusController struct {
	controller.HTTP
	processor processors.ChangeSubmissionStatusProcessorInterface
}

func (C ChangeSubmissionStatusController) New(processor processors.ChangeSubmissionStatusProcessorInterface) *ChangeSubmissionStatusController {
	c := new(ChangeSubmissionStatusController)
	c.processor = processor
	return c
}

func (C *ChangeSubmissionStatusController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.ChangeSubmissionStatus(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *ChangeSubmissionStatusController) ChangeSubmissionStatus(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	submissionID := fiberhelper.Param(c, "id")
	if submissionID == "" {
		return controller.ErrorResponse(fiber.StatusBadRequest, "id not provided")
	}

	submissionStatusRequest := new(requests.SubmissionStatusRequest)
	err = fiberhelper.BindJSON(c, submissionStatusRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	submission, err := C.processor.ChangeSubmissionStatus(player, submissionID, submissionStatusRequest.Status)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Submission *models.SubmissionView
	}{
		submission,
	})
}
