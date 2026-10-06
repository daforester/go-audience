package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ToggleBadgeCheckController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C ToggleBadgeCheckController) New(gameRepository repositories.GameRepository) *ToggleBadgeCheckController {
	c := new(ToggleBadgeCheckController)
	c.gameRepository = gameRepository
	return c
}

func (C *ToggleBadgeCheckController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ToggleBadgeCheck(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ToggleBadgeCheckController) ToggleBadgeCheck(c *gin.Context) *controller.Response {
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

	toggleRequest := &struct {
		Enabled bool
	}{
		true,
	}
	err := c.BindJSON(toggleRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	C.gameRepository.Lock()
	game.BadgeCheck = toggleRequest.Enabled
	C.gameRepository.Save(game)
	C.gameRepository.Unlock()

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
