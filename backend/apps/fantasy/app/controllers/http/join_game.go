package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/requests"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
	"strings"
)

type JoinGameController struct {
	controller.HTTP
	processor      processors.JoinGameProcessorInterface
	loginProcessor processors.LoginProcessorInterface
}

func (C JoinGameController) New(processor processors.JoinGameProcessorInterface, loginProcessor processors.LoginProcessorInterface) *JoinGameController {
	c := new(JoinGameController)
	c.processor = processor
	c.loginProcessor = loginProcessor
	return c
}

func (C *JoinGameController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.JoinGame(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *JoinGameController) JoinGame(c fiber.Ctx) *controller.Response {
	joinGameRequest := new(requests.JoinGameRequest)
	err := fiberhelper.BindJSON(c, joinGameRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	joinGameRequest.RoomCode = strings.ToUpper(strings.TrimSpace(joinGameRequest.RoomCode))

	player, err := C.processor.JoinGame(joinGameRequest.Name, joinGameRequest.RoomCode)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	claims := make(map[string]interface{})
	claims["RoomKey"] = player.RoomKey
	tokens, err := C.loginProcessor.LoginTokens(player, player.UUID, true, claims)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Player *models.Player
		Tokens *processors.Tokens
	}{
		player,
		tokens,
	})
}
