package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/requests"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
)

type SetPresentingTeamController struct {
	controller.HTTP
	processor processors.SetPresentingTeamProcessorInterface
}

func (C SetPresentingTeamController) New(processor processors.SetPresentingTeamProcessorInterface) *SetPresentingTeamController {
	c := new(SetPresentingTeamController)
	c.processor = processor
	return c
}

func (C *SetPresentingTeamController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.SetPresentingTeam(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *SetPresentingTeamController) SetPresentingTeam(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	presentingTeamRequest := new(requests.PresentingTeamRequest)
	err = fiberhelper.BindJSON(c, presentingTeamRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	err = C.processor.SetPresentingTeam(player, presentingTeamRequest.TeamID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		TeamID string
	}{
		presentingTeamRequest.TeamID,
	})
}
