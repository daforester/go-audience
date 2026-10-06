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
	c := C.GetContext().(*gin.Context)
	e := C.JoinGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *JoinGameController) JoinGame(c *gin.Context) *controller.Response {
	joinGameRequest := new(requests.JoinGameRequest)
	err := c.BindJSON(joinGameRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}
	joinGameRequest.Name = strings.TrimSpace(joinGameRequest.Name)
	joinGameRequest.RoomCode = strings.ToUpper(strings.TrimSpace(joinGameRequest.RoomCode))

	player, err := C.processor.JoinGame(joinGameRequest.Name, joinGameRequest.RoomCode)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	claims := make(map[string]interface{})
	claims["RoomKey"] = joinGameRequest.RoomCode
	tokens, err := C.loginProcessor.LoginTokens(player, player.UUID, true, claims)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Player *models.Player
		Tokens *processors.Tokens
	}{
		player,
		tokens,
	})
}
