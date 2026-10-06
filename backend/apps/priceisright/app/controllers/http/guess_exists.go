package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GuessExistsController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C GuessExistsController) New(gameRepository repositories.GameRepository) *GuessExistsController {
	c := new(GuessExistsController)
	c.gameRepository = gameRepository
	return c
}

func (C *GuessExistsController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.GuessExists(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *GuessExistsController) GuessExists(c *gin.Context) *controller.Response {
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

	game := C.gameRepository.Get()

	if game == nil {
		return controller.ErrorResponse(http.StatusUnauthorized, "game not found")
	}

	round := game.CurrentRound
	if round == nil {
		return controller.StandardResponse(struct {
			Guess *models.Guess
		}{
			nil,
		})
	}

	for _, guess := range round.Guesses {
		if guess.PlayerUUID == player.UUID {
			return controller.StandardResponse(struct {
				Guess *models.Guess
			}{
				guess,
			})
		}
	}

	return controller.StandardResponse(struct {
		Guess *models.Guess
	}{
		nil,
	})
}
