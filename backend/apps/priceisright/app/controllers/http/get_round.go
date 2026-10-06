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

type GetRoundController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	processor      processors.SetRoundProcessor
}

func (C GetRoundController) New(gameRepository repositories.GameRepository, processor processors.SetRoundProcessor) *GetRoundController {
	c := new(GetRoundController)
	c.gameRepository = gameRepository
	c.processor = processor
	return c
}

func (C *GetRoundController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.GetRound(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *GetRoundController) GetRound(c *gin.Context) *controller.Response {
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

	if game.RoundNumber > len(game.Rounds)-1 {
		if len(game.Rounds) > 0 {
			_, err := C.processor.SetRound(0)
			if err != nil {
				return controller.ErrorResponse(http.StatusUnauthorized, "unable to set round to zero")
			}
		} else {
			return controller.ErrorResponse(http.StatusUnauthorized, "no available round")
		}
	}

	if player.Class == models.PLAYER_CLASS_LEADER || player.Class == models.PLAYER_CLASS_DISPLAY {
		return controller.StandardResponse(struct {
			Round *models.Round
		}{
			game.Rounds[game.RoundNumber],
		})
	}

	round := &models.Round{
		Image: game.Rounds[game.RoundNumber].Image,
		Name:  game.Rounds[game.RoundNumber].Name,
		UUID:  game.Rounds[game.RoundNumber].UUID,
	}

	return controller.StandardResponse(struct {
		Round *models.Round
	}{
		round,
	})
}
