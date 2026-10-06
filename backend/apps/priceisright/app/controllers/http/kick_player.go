package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type KickPlayerController struct {
	controller.HTTP
	kickProcessor processors.KickPlayerProcessorInterface
	quitProcessor processors.QuitGameProcessorInterface
}

func (C KickPlayerController) New(kickProcessor processors.KickPlayerProcessorInterface, quitProcessor processors.QuitGameProcessorInterface) *KickPlayerController {
	c := new(KickPlayerController)
	c.kickProcessor = kickProcessor
	c.quitProcessor = quitProcessor
	return c
}

func (C *KickPlayerController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.KickPlayer(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *KickPlayerController) KickPlayer(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be room leader to kick players")
	}

	uuid := c.Param("uuid")
	if uuid == "" {
		return controller.ErrorResponse(http.StatusUnauthorized, "uuid not provided")
	}

	if player.UUID == uuid {
		fmt.Println("Quit Game Request")
		err := C.quitProcessor.QuitGame(player)
		if err != nil {
			return controller.ErrorResponse(http.StatusBadRequest, err.Error())
		}
	} else {
		fmt.Println("Kick Player Request")
		err := C.kickProcessor.KickPlayer(player, uuid)
		if err != nil {
			return controller.ErrorResponse(http.StatusBadRequest, err.Error())
		}
	}

	return controller.StandardResponse(struct {
		Removed bool
	}{
		true,
	})
}
