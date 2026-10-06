package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
)

type RemoveTeamController struct {
	controller.HTTP
	processor processors.TeamProcessorInterface
}

func (C RemoveTeamController) New(processor processors.TeamProcessorInterface) *RemoveTeamController {
	c := new(RemoveTeamController)
	c.processor = processor
	return c
}

func (C *RemoveTeamController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.RemoveTeam(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *RemoveTeamController) RemoveTeam(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	teamID := fiberhelper.Param(c, "id")
	if teamID == "" {
		return controller.ErrorResponse(fiber.StatusBadRequest, "id not provided")
	}

	game, err := C.processor.RemoveTeam(player, teamID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
