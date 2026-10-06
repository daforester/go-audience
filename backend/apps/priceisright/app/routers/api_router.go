package routers

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/controllers/http"
	"bitbucket.org/daforester/go-audience/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-gonic/gin"
)

type APIRouter struct {
	router.BaseRouter
}

func NewAPIRouter(app di.AppInterface, engine ...*gin.Engine) router.Router {
	r := (&APIRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P *APIRouter) New(app di.AppInterface) router.Router {
	p := new(APIRouter)
	p.SetContainer(app)
	return p
}

func (P *APIRouter) RegisterRoutes(router *gin.Engine) {
	P.RegisterPublicRoutes(router)
}

func (P *APIRouter) RegisterPublicRoutes(router *gin.Engine) {
	router.GET("/api/token/refresh", func(c *gin.Context) {
		P.RunController(c, &http.RefreshTokenController{})
	})

	router.GET("/api/game", func(c *gin.Context) {
		P.RunController(c, &http.GetGameController{})
	})

	router.POST("/api/game/code", func(c *gin.Context) {
		P.RunController(c, &http.SetGameCodeController{})
	})

	router.POST("/api/game/join", func(c *gin.Context) {
		P.RunController(c, &http.JoinGameController{})
	})

	router.GET("/api/game/rejoin", func(c *gin.Context) {
		P.RunController(c, &http.RejoinGameController{})
	})

	router.GET("/api/game/quit", func(c *gin.Context) {
		P.RunController(c, &http.QuitGameController{})
	})

	router.GET("/api/game/kick/:uuid", func(c *gin.Context) {
		P.RunController(c, &http.KickPlayerController{})
	})

	router.POST("/api/game/player/change/:index", func(c *gin.Context) {
		P.RunController(c, &http.ChangePlayerNameController{})
	})

	router.POST("/api/game/player/badge/change/:index", func(c *gin.Context) {
		P.RunController(c, &http.ChangePlayerBadgeController{})
	})

	router.POST("/api/game/password", func(c *gin.Context) {
		P.RunController(c, &http.ChangeGamePasswordController{})
	})

	router.POST("/api/game/role/assign", func(c *gin.Context) {
		P.RunController(c, &http.AssignRoleController{})
	})

	router.POST("/api/game/rounds/save", func(c *gin.Context) {
		P.RunController(c, &http.SaveRoundsController{})
	})

	router.POST("/api/game/guess/change/:index", func(c *gin.Context) {
		P.RunController(c, &http.ChangeGuessController{})
	})

	router.PUT("/api/game/guess/new", func(c *gin.Context) {
		P.RunController(c, &http.NewGuessController{})
	})

	router.GET("/api/game/guess/exists", func(c *gin.Context) {
		P.RunController(c, &http.GuessExistsController{})
	})

	router.POST("/api/game/guesses/clear", func(c *gin.Context) {
		P.RunController(c, &http.ClearGuessesController{})
	})

	router.GET("/api/game/round", func(c *gin.Context) {
		P.RunController(c, &http.GetRoundController{})
	})

	router.POST("/api/game/round/set", func(c *gin.Context) {
		P.RunController(c, &http.SetRoundController{})
	})

	router.POST("/api/game/round/activate", func(c *gin.Context) {
		P.RunController(c, &http.ActivateRoundController{})
	})

	router.GET("/api/game/round/guesses/:round", func(c *gin.Context) {
		P.RunController(c, &http.TopGuessesController{})
	})

	router.POST("/api/game/summon", func(c *gin.Context) {
		P.RunController(c, &http.SummonPlayerController{})
	})

	router.POST("/api/game/badgecheck", func(c *gin.Context) {
		P.RunController(c, &http.ToggleBadgeCheckController{})
	})
}
