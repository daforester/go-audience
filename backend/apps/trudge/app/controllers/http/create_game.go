package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/requests"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	c := C.GetContext().(*gin.Context)
	e := C.CreateGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *CreateGameController) CreateGame(c *gin.Context) *controller.Response {
	createGameRequest := new(requests.CreateGameRequest)
	err := c.BindJSON(createGameRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	createGameRequest.Name = strings.TrimSpace(createGameRequest.Name)
	game, player, err := C.processor.CreateGame(createGameRequest.Name)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	claims := make(map[string]interface{})
	claims["RoomKey"] = game.Key
	tokens, err := C.loginProcessor.LoginTokens(player, player.UUID, true, claims)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game   *models.Game
		Player *models.Player
		Tokens *processors.Tokens
	}{
		game,
		player,
		tokens,
	})
}
