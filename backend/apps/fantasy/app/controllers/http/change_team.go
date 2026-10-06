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

type ChangeTeamController struct {
	controller.HTTP
	processor processors.TeamProcessorInterface
}

func (C ChangeTeamController) New(processor processors.TeamProcessorInterface) *ChangeTeamController {
	c := new(ChangeTeamController)
	c.processor = processor
	return c
}

func (C *ChangeTeamController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.ChangeTeam(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *ChangeTeamController) ChangeTeam(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	teamID := fiberhelper.Param(c, "id")
	if teamID == "" {
		return controller.ErrorResponse(fiber.StatusBadRequest, "id not provided")
	}

	teamRequest := new(requests.TeamRequest)
	err = fiberhelper.BindJSON(c, teamRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeTeam(player, teamID, teamRequest.Name)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
