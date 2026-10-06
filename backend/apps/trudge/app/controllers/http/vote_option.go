package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VoteOptionController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	processor      processors.VoteOptionProcessorInterface
}

func (C VoteOptionController) New(gameRepository repositories.GameRepository, processor processors.VoteOptionProcessorInterface) *VoteOptionController {
	c := new(VoteOptionController)
	c.gameRepository = gameRepository
	c.processor = processor
	return c
}

func (C *VoteOptionController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.VoteOption(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *VoteOptionController) VoteOption(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	voteOptionRequest := &struct {
		Index  int
		Option string
	}{
		0,
		"",
	}
	err = c.BindJSON(voteOptionRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, err = C.processor.VoteOption(player.RoomKey, player, voteOptionRequest.Index, voteOptionRequest.Option)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Option string
	}{
		voteOptionRequest.Option,
	})
}
