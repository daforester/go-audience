package router

import (
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-gonic/gin"
	netHttp "net/http"
)

type Router interface {
	RegisterRoutes(*gin.Engine)
	SetContainer(container di.AppInterface)
	RunController(*gin.Context, controller.Controller)
}

type BaseRouter struct {
	container di.AppInterface
}

func (B *BaseRouter) RegisterRoutes(engine *gin.Engine) {

}

func (B *BaseRouter) SetContainer(container di.AppInterface) {
	B.container = container
}

func (B *BaseRouter) RunController(gin *gin.Context, C controller.Controller) {
	c := B.container.Make(C)
	c.(controller.Controller).SetContext(gin)

	err := c.(controller.Controller).Run()
	if err != nil {
		gin.JSON(netHttp.StatusInternalServerError, err.Error())
	}
}
