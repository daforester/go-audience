package routers

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/controllers/http"
	"bitbucket.org/daforester/go-audience/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-gonic/gin"
)

type WebSocketRouter struct {
	router.BaseRouter
}

func NewWebSocketRouter(app di.AppInterface, engine ...*gin.Engine) router.Router {
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

func (P *WebSocketRouter) RegisterRoutes(router *gin.Engine) {
	P.RegisterPublicRoutes(router)
}

func (P *WebSocketRouter) RegisterPublicRoutes(router *gin.Engine) {
	router.GET("/ws/:room/:player", func(c *gin.Context) {
		P.RunController(c, &http.EstablishSocketController{})
	})
}
