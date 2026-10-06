package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
)

type StaffVoteController struct {
	controller.HTTP
	processor processors.StaffVoteProcessorInterface
}

func (C StaffVoteController) New(processor processors.StaffVoteProcessorInterface) *StaffVoteController {
	c := new(StaffVoteController)
	c.processor = processor
	return c
}

func (C *StaffVoteController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.StaffVote(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *StaffVoteController) StaffVote(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	submissionID := fiberhelper.Param(c, "id")
	if submissionID == "" {
		return controller.ErrorResponse(fiber.StatusBadRequest, "id not provided")
	}

	submission, err := C.processor.StaffVote(player, submissionID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Submission *models.SubmissionView
	}{
		submission,
	})
}
