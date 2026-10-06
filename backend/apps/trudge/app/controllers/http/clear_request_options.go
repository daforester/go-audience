package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	processors "bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	models "bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"bitbucket.org/daforester/go-audience/services/ginhelper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClearRequestOptionsController struct {
	controller.HTTP
	processor     processors.ClearRequestOptionsProcessorInterface
	codeProcessor processors.SetGameCodeProcessorInterface
}

func (C ClearRequestOptionsController) New(processor processors.ClearRequestOptionsProcessorInterface, codeProcessor processors.SetGameCodeProcessorInterface) *ClearRequestOptionsController {
	c := new(ClearRequestOptionsController)
	c.processor = processor
	c.codeProcessor = codeProcessor

	return c
}

func (C *ClearRequestOptionsController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.ShowRequestOptions(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *ClearRequestOptionsController) ShowRequestOptions(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	prompt := ""
	data, _ := ginhelper.GetRequestBody(c)
	p, e := data["prompt"]

	if e {
		prompt = p.(string)
	}

	game, err := C.processor.ClearRequestOptions(player.RoomKey, prompt)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
