package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChangeQuestController struct {
	controller.HTTP
	processor processors.ChangeQuestProcessorInterface
}

func (C ChangeQuestController) New(processor processors.ChangeQuestProcessorInterface) *ChangeQuestController {
	c := new(ChangeQuestController)
	c.processor = processor

	return c
}

func (C *ChangeQuestController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ChangeQuest(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ChangeQuestController) ChangeQuest(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	ChangeQuestRequest := &struct {
		Quest string
	}{
		"",
	}
	err = c.BindJSON(ChangeQuestRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.ChangeQuest(player.RoomKey, ChangeQuestRequest.Quest)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
