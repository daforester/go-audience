package routers

import (
	"bitbucket.org/daforester/go-audience/services/ginhelper"
	"bitbucket.org/daforester/go-audience/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
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
	router.GET("/", func(c *gin.Context) {
		dashboardPage(c)
	})
	router.GET("/index.html", func(c *gin.Context) {
		dashboardPage(c)
	})
	router.GET("/index.gohtml", func(c *gin.Context) {
		dashboardPage(c)
	})
	router.GET("/mapmaker", func(c *gin.Context) {
		dashboardPage(c)
	})
	router.GET("/advert", func(c *gin.Context) {
		dashboardPage(c)
	})
	router.GET("/logout", func(c *gin.Context) {
		dashboardPage(c)
	})

	router.GET("/service-worker.js", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		c.File(viper.GetString("ROOT_PATH") + "/public/service-worker.js")
	})

	router.GET("/js/trudge/:filename", func(c *gin.Context) {
		jsClientPath := viper.GetString("ROOT_PATH") + "/public/js/trudge/"
		mainJs := "main.[a-z0-9]+.js"
		fileName := c.Param("filename")
		if fileName == "main.js" {
			filepath.Walk(jsClientPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal(err)
					return err
				}
				match, err := regexp.Match(mainJs, []byte(filepath.Base(path)))
				if err != nil {
					log.Fatal(err)
					return err
				}
				if match {
					fileName = filepath.Base(path)
				}
				return nil
			})
		}

		modified, eTag := ginhelper.FileModified(c, jsClientPath+fileName)
		if modified {
			c.Header("Cache-Control", "no-cache")
			c.Header("Etag", "\""+eTag+"\"")
			c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(fileName)))
			c.File(jsClientPath + fileName)
		} else {
			c.Status(http.StatusNotModified)
		}
	})

	router.GET("/fonts/:filename", func(c *gin.Context) {
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/fonts/" + c.Param("filename"))
	})

	router.GET("/images/app/:filename", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=31536000")
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/images/app/" + c.Param("filename"))
	})

	router.GET("/images/static/:filename", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=31536000")
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/images/static/" + c.Param("filename"))
	})

	router.GET("/images/maps/:filename", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=31536000")
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/images/maps/" + c.Param("filename"))
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=86400")
		c.Header("Content-Type", mime.TypeByExtension(".ico"))
		c.File(viper.GetString("ROOT_PATH") + "/public/favicon.ico")
	})

	router.GET("/music/:filename", func(c *gin.Context) {
		c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(c.Param("filename"))))
		c.File(viper.GetString("ROOT_PATH") + "/public/music/" + c.Param("filename"))
	})

	router.NoRoute(func(c *gin.Context) {
		r, _ := regexp.Compile("^/manifest\\.[a-z0-9]+\\.json$")
		if r != nil && r.Match([]byte(c.Request.RequestURI)) {
			c.Header("Cache-Control", "max-age=31536000")
			c.Header("Content-Type", mime.TypeByExtension(".json"))
			c.File(viper.GetString("ROOT_PATH") + "/public/" + c.Request.RequestURI)
			return
		}

		r, _ = regexp.Compile("^/precache-manifest\\.[a-z0-9]+\\.js$")
		if r != nil && r.Match([]byte(c.Request.RequestURI)) {
			c.Header("Cache-Control", "max-age=31536000")
			c.Header("Content-Type", mime.TypeByExtension(".js"))
			c.File(viper.GetString("ROOT_PATH") + "/public/" + c.Request.RequestURI)
			return
		}

		r, _ = regexp.Compile("^/[a-z0-9]+\\.woff?$")
		if r != nil && r.Match([]byte(c.Request.RequestURI)) {
			c.Header("Cache-Control", "max-age=31536000")
			c.Header("Content-Type", mime.TypeByExtension(".woff"))
			c.File(viper.GetString("ROOT_PATH") + "/public/" + c.Request.RequestURI)
			return
		}

		r, _ = regexp.Compile("^/[a-z0-9]+\\.woff2?$")
		if r != nil && r.Match([]byte(c.Request.RequestURI)) {
			c.Header("Cache-Control", "max-age=31536000")
			c.Header("Content-Type", mime.TypeByExtension(".woff2"))
			c.File(viper.GetString("ROOT_PATH") + "/public/" + c.Request.RequestURI)
			return
		}

		c.HTML(http.StatusNotFound, "index.gohtml", map[string]string{"XErrorMessage": "Page not found"})
	})
}

func dashboardPage(c *gin.Context) {
	data := make(map[string]interface{})
	data["PageTitle"] = "S.L.O.S.H. Members"
	data["Year"] = time.Now().Local().Format("2006")
	data["API_URI"] = viper.GetString("API_URL") + "/api"
	if viper.GetString("CLIENT_TOKENS") != "" {
		data["ClientTokens"] = template.JS(viper.GetString("CLIENT_TOKENS"))
	}

	c.HTML(http.StatusOK, "index.gohtml", data)
}
