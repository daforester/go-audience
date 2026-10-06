package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
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
	c := C.GetContext().(fiber.Ctx)
	e := C.KickPlayer(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *KickPlayerController) KickPlayer(c fiber.Ctx) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	uuid := fiberhelper.Param(c, "uuid")
	if uuid == "" {
		return controller.ErrorResponse(fiber.StatusBadRequest, "uuid not provided")
	}

	if player.UUID == uuid {
		err = C.quitProcessor.QuitGame(player)
	} else {
		err = C.kickProcessor.KickPlayer(player, uuid)
	}

	if err != nil {
		return helpers.ErrorResponse(err)
	}

	return controller.StandardResponse(struct {
		Removed bool
	}{
		true,
	})
}
