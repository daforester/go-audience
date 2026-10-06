package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SyncVolumeController struct {
	controller.HTTP
	processor processors.SyncVolumeProcessorInterface
}

func (C SyncVolumeController) New(processor processors.SyncVolumeProcessorInterface) *SyncVolumeController {
	c := new(SyncVolumeController)
	c.processor = processor
	return c
}

func (C *SyncVolumeController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.SyncVolume(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *SyncVolumeController) SyncVolume(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_DISPLAY &&
		player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be a display or game leader")
	}

	setVolumeRequest := &struct {
		Volume int
	}{
		0,
	}
	err = c.BindJSON(setVolumeRequest)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	game, err := C.processor.SyncVolume(player.RoomKey, setVolumeRequest.Volume)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return controller.StandardResponse(struct {
		Game *models.Game
	}{
		game,
	})
}
