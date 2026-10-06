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

type CreateGameController struct {
	controller.HTTP
	loginProcessor processors.LoginProcessorInterface
	processor      processors.CreateGameProcessorInterface
}

func (C CreateGameController) New(processor processors.CreateGameProcessorInterface, loginProcessor processors.LoginProcessorInterface) *CreateGameController {
	c := new(CreateGameController)
	c.processor = processor
	c.loginProcessor = loginProcessor
	return c
}

func (C *CreateGameController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.CreateGame(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *CreateGameController) CreateGame(c fiber.Ctx) *controller.Response {
	createGameRequest := new(requests.CreateGameRequest)
	err := fiberhelper.BindJSON(c, createGameRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, player, err := C.processor.CreateGame(createGameRequest.Name)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	claims := make(map[string]interface{})
	claims["RoomKey"] = game.Key
	tokens, err := C.loginProcessor.LoginTokens(player, player.UUID, true, claims)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game   *models.GameView
		Player *models.Player
		Tokens *processors.Tokens
	}{
		game,
		player,
		tokens,
	})
}
