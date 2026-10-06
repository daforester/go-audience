package router

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"github.com/daforester/go-di-container/di"
	"github.com/gofiber/fiber/v3"
)

type Router interface {
	RegisterRoutes(*fiber.App)
	SetContainer(container di.AppInterface)
	RunController(fiber.Ctx, controller.Controller) error
}

type BaseRouter struct {
	container di.AppInterface
}

func (B *BaseRouter) RegisterRoutes(app *fiber.App) {

}

func (B *BaseRouter) SetContainer(container di.AppInterface) {
	B.container = container
}

func (B *BaseRouter) RunController(ctx fiber.Ctx, C controller.Controller) error {
	c := B.container.Make(C)
	c.(controller.Controller).SetContext(ctx)

	err := c.(controller.Controller).Run()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return nil
}
