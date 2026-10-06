package routers

import (
	"bitbucket.org/daforester/go-audience/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mime"
	"path/filepath"
)

type ViewRouter struct {
	router.BaseRouter
}

func NewViewRouter(app di.AppInterface, engine ...*gin.Engine) router.Router {
	r := (&ViewRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P *ViewRouter) New(app di.AppInterface) router.Router {
	a := new(ViewRouter)
	a.SetContainer(app)
	return a
}

func (P *ViewRouter) RegisterRoutes(router *gin.Engine) {
	router.GET("/images/rounds/:filename", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=31536000")
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/images/rounds/" + c.Param("filename"))
	})
}
