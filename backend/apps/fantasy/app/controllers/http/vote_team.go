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

type VoteTeamController struct {
	controller.HTTP
	processor processors.VoteTeamProcessorInterface
}

func (C VoteTeamController) New(processor processors.VoteTeamProcessorInterface) *VoteTeamController {
	c := new(VoteTeamController)
	c.processor = processor
	return c
}

func (C *VoteTeamController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.VoteTeam(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *VoteTeamController) VoteTeam(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	voteTeamRequest := new(requests.VoteTeamRequest)
	err = fiberhelper.BindJSON(c, voteTeamRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	round, err := C.processor.VoteTeam(player, voteTeamRequest.TeamID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Round *models.RoundView
	}{
		round,
	})
}
