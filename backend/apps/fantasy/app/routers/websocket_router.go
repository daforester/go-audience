package routers

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/controllers/http"
	"bitbucket.org/daforester/go-audience/apps/fantasy/middleware"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gofiber/fiber/v3"
)

type WebSocketRouter struct {
	router.BaseRouter
}

func NewWebSocketRouter(app di.AppInterface, engine ...*fiber.App) router.Router {
	r := (&WebSocketRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P *WebSocketRouter) New(app di.AppInterface) router.Router {
	p := new(WebSocketRouter)
	p.SetContainer(app)
	return p
}

func (P *WebSocketRouter) RegisterRoutes(router *fiber.App) {
	P.RegisterPublicRoutes(router)
}

func (P *WebSocketRouter) RegisterPublicRoutes(router *fiber.App) {
	router.Get("/ws/:room/:player", middleware.WebSocketUpgrade(), func(c fiber.Ctx) error {
		return P.RunController(c, &http.EstablishSocketController{})
	})
}
