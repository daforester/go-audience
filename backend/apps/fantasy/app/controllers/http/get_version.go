package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"github.com/gofiber/fiber/v3"
)

type GetVersionController struct {
	controller.HTTP
}

func (C GetVersionController) New() *GetVersionController {
	return new(GetVersionController)
}

func (C *GetVersionController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.GetVersion()

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *GetVersionController) GetVersion() *controller.Response {
	return controller.StandardResponse(struct {
		Version string
	}{
		setup.Version,
	})
}
