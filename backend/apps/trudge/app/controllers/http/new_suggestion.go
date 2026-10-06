package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewSuggestionController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	processor      processors.NewSuggestionProcessorInterface
}

func (C NewSuggestionController) New(gameRepository repositories.GameRepository, processor processors.NewSuggestionProcessorInterface) *NewSuggestionController {
	c := new(NewSuggestionController)
	c.gameRepository = gameRepository
	c.processor = processor
	return c
}

func (C *NewSuggestionController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.NewSuggestion(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *NewSuggestionController) NewSuggestion(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_AUDIENCE {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be an audience member")
	}

	NewSuggestionRequest := &struct {
		Suggestion string
	}{
		"",
	}
	err = c.BindJSON(NewSuggestionRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, err = C.processor.NewSuggestion(player.RoomKey, NewSuggestionRequest.Suggestion)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Suggestion string
	}{
		NewSuggestionRequest.Suggestion,
	})
}
