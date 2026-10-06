package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/middleware"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RejoinGameController struct {
	controller.HTTP
	repository repositories.GameRepository
}

func (C RejoinGameController) New(repository repositories.GameRepository) *RejoinGameController {
	c := new(RejoinGameController)
	c.repository = repository
	return c
}

func (C *RejoinGameController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.RejoinGame(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *RejoinGameController) RejoinGame(c *gin.Context) *controller.Response {
	playerInterface, exists := c.Get("Player")
	if !exists {
		err := middleware.ProcessPlayerToken(c)
		if err != nil {
			return controller.ErrorResponse(http.StatusUnauthorized, "player token data not provided")
		}
		playerInterface, exists = c.Get("Player")
		if !exists {
			return controller.ErrorResponse(http.StatusUnauthorized, "player token data not parsed")
		}
	}

	player := playerInterface.(*models.Player)

	game := C.repository.GetByKey(player.RoomKey)

	if game == nil {
		return controller.ErrorResponse(http.StatusUnauthorized, "game does not exist")
	}

	for _, p := range game.Players {
		if p.UUID == player.UUID {
			player = p
		}
	}

	if player.Class == models.PLAYER_CLASS_LEADER || player.Class == models.PLAYER_CLASS_DISPLAY {
		return controller.StandardResponse(struct {
			Game   *models.Game
			Player *models.Player
		}{
			game,
			player,
		})
	}

	return controller.StandardResponse(struct {
		GameCode           int
		Player             *models.Player
		SuggestionsEnabled bool
	}{
		game.GameCode,
		player,
		game.SuggestionsEnabled,
	})
}
