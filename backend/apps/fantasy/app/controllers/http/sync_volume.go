package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/requests"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
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
	c := C.GetContext().(fiber.Ctx)
	e := C.SyncVolume(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *SyncVolumeController) SyncVolume(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	syncVolumeRequest := new(requests.SyncVolumeRequest)
	err = fiberhelper.BindJSON(c, syncVolumeRequest)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}

	game, err := C.processor.SyncVolume(player, syncVolumeRequest.Volume)
	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Game *models.GameView
	}{
		game,
	})
}
