package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TopGuessesController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C TopGuessesController) New(gameRepository repositories.GameRepository) *TopGuessesController {
	c := new(TopGuessesController)
	c.gameRepository = gameRepository
	return c
}

func (C *TopGuessesController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.TopGuesses(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *TopGuessesController) TopGuesses(c *gin.Context) *controller.Response {
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

	round := game.RoundNumber

	roundNumber := c.Param("round")
	if roundNumber != "" {
		round, _ = strconv.Atoi(roundNumber)
	}

	if round < 0 || round >= len(game.Rounds) {
		return controller.ErrorResponse(http.StatusUnauthorized, "round out of range")
	}

	if player.Class == models.PLAYER_CLASS_LEADER || player.Class == models.PLAYER_CLASS_DISPLAY {
		return controller.StandardResponse(struct {
			Guesses []*models.Guess
		}{
			game.Rounds[round].TopGuesses(),
		})
	}

	return controller.ErrorResponse(http.StatusUnauthorized, "access denied")
}
