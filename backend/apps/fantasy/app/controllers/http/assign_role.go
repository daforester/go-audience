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

type AssignRoleController struct {
	controller.HTTP
	processor processors.AssignRoleProcessorInterface
}

func (C AssignRoleController) New(processor processors.AssignRoleProcessorInterface) *AssignRoleController {
	c := new(AssignRoleController)
	c.processor = processor
	return c
}

func (C *AssignRoleController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.AssignRole(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *AssignRoleController) AssignRole(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	assignPlayerRequest := new(requests.AssignPlayerRequest)
	err = fiberhelper.BindJSON(c, assignPlayerRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	p, err := C.processor.AssignPlayer(player, assignPlayerRequest.UUID, assignPlayerRequest.PlayerClass, assignPlayerRequest.TeamID)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Player *models.Player
	}{
		p,
	})
}
