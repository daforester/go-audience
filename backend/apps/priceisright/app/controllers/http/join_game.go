package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/requests"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JoinGameController struct {
	controller.HTTP
	processor      processors.JoinGameProcessorInterface
	loginProcessor processors.LoginProcessorInterface
	gameRepository repositories.GameRepository
}

func (C JoinGameController) New(processor processors.JoinGameProcessorInterface, loginProcessor processors.LoginProcessorInterface, gameRepository repositories.GameRepository) *JoinGameController {
	c := new(JoinGameController)
	c.processor = processor
	c.loginProcessor = loginProcessor
	c.gameRepository = gameRepository
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
	if joinGameRequest.BadgeNumber == "" {
		joinGameRequest.BadgeNumber = "000X"
	}
	joinGameRequest.Name = strings.TrimSpace(joinGameRequest.Name)
	joinGameRequest.BadgeNumber = strings.ToUpper(strings.TrimSpace(joinGameRequest.BadgeNumber))
	joinGameRequest.Password = strings.TrimSpace(joinGameRequest.Password)

	player, err := C.processor.JoinGame(joinGameRequest.Name, joinGameRequest.BadgeNumber, joinGameRequest.Password)
	if err != nil {
		if err.Error() == "INVALID_BADGE" {
			return controller.ErrorResponse(http.StatusNotAcceptable, "Badge not found")
		}
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	claims := make(map[string]interface{})
	claims["BadgeNumber"] = joinGameRequest.BadgeNumber
	tokens, err := C.loginProcessor.LoginTokens(player, player.UUID, true, claims)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if player.Class == models.PLAYER_CLASS_LEADER {
		game := C.gameRepository.Get()
		return controller.StandardResponse(struct {
			Game   *models.Game
			Player *models.Player
			Tokens *processors.Tokens
		}{
			game,
			player,
			tokens,
		})
	} else {
		return controller.StandardResponse(struct {
			Player *models.Player
			Tokens *processors.Tokens
		}{
			player,
			tokens,
		})
	}
}
