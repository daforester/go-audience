package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClearGuessesController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	processor      processors.ClearGuessesProcessorInterface
}

func (C ClearGuessesController) New(gameRepository repositories.GameRepository, processorInterface processors.ClearGuessesProcessorInterface) *ClearGuessesController {
	c := new(ClearGuessesController)
	c.gameRepository = gameRepository
	c.processor = processorInterface
	return c
}

func (C *ClearGuessesController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ClearGuesses(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ClearGuessesController) ClearGuesses(c *gin.Context) *controller.Response {
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

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "access denied")
	}

	game := C.gameRepository.Get()

	if game == nil {
		return controller.ErrorResponse(http.StatusUnauthorized, "game not found")
	}

	guess, err := C.processor.ClearGuesses(game.RoundNumber)
	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, "problem clearing guesses")
	}

	return controller.StandardResponse(struct {
		Guess []*models.Guess
	}{
		guess,
	})
}
