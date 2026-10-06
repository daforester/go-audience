package controller

import (
	"github.com/gofiber/fiber/v3"
)

type HTTP struct {
	BaseController
	c       fiber.Ctx
	context fiber.Ctx
}

func (C *HTTP) GetContext() interface{} {
	return C.context
}

func (C *HTTP) SetContext(context interface{}) Controller {
	C.c = context.(fiber.Ctx)
	C.context = context.(fiber.Ctx)
	return C
}

func (C *HTTP) Run() error {
	return nil
}
