package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ToggleSuggestionsController struct {
	controller.HTTP
	processor processors.ToggleSuggestionsProcessorInterface
}

func (C ToggleSuggestionsController) New(processor processors.ToggleSuggestionsProcessorInterface) *ToggleSuggestionsController {
	c := new(ToggleSuggestionsController)
	c.processor = processor
	return c
}

func (C *ToggleSuggestionsController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ToggleSuggestions(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ToggleSuggestionsController) ToggleSuggestions(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_DISPLAY &&
		player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be a display or game leader")
	}

	toggleSuggestionsRequest := &struct {
		Enabled bool
	}{
		false,
	}
	err = c.BindJSON(toggleSuggestionsRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ToggleSuggestions(player.RoomKey, toggleSuggestionsRequest.Enabled)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
