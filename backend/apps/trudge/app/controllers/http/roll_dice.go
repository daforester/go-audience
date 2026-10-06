package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RollDiceController struct {
	controller.HTTP
	processor processors.DiceRollProcessorInterface
}

func (C RollDiceController) New(processor processors.DiceRollProcessorInterface) *RollDiceController {
	c := new(RollDiceController)
	c.processor = processor
	return c
}

func (C *RollDiceController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.RollDice(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *RollDiceController) RollDice(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_BARD &&
		player.Class != models.PLAYER_CLASS_WARRIOR &&
		player.Class != models.PLAYER_CLASS_WIZARD &&
		player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be a party member or game leader")
	}

	game, err := C.processor.RollDice(player.RoomKey)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
