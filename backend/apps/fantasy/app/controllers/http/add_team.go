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

type AddTeamController struct {
	controller.HTTP
	processor processors.TeamProcessorInterface
}

func (C AddTeamController) New(processor processors.TeamProcessorInterface) *AddTeamController {
	c := new(AddTeamController)
	c.processor = processor
	return c
}

func (C *AddTeamController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.AddTeam(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *AddTeamController) AddTeam(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	teamRequest := new(requests.TeamRequest)
	err = fiberhelper.BindJSON(c, teamRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.AddTeam(player, teamRequest.Name)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
