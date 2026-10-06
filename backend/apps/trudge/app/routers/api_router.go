package routers

import (
	http2 "bitbucket.org/daforester/go-audience/apps/trudge/app/controllers/http"
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
		P.RunController(c, &http2.RefreshTokenController{})
	})

	router.GET("/api/game", func(c *gin.Context) {
		P.RunController(c, &http2.GetGameController{})
	})

	router.POST("/api/game/create", func(c *gin.Context) {
		P.RunController(c, &http2.CreateGameController{})
	})

	router.GET("/api/game/start", func(c *gin.Context) {
		P.RunController(c, &http2.StartGameController{})
	})

	router.POST("/api/game/code", func(c *gin.Context) {
		P.RunController(c, &http2.SetGameCodeController{})
	})

	router.GET("/api/game/kick/:uuid", func(c *gin.Context) {
		P.RunController(c, &http2.KickPlayerController{})
	})

	router.POST("/api/game/join", func(c *gin.Context) {
		P.RunController(c, &http2.JoinGameController{})
	})

	router.GET("/api/game/rejoin", func(c *gin.Context) {
		P.RunController(c, &http2.RejoinGameController{})
	})

	router.GET("/api/game/quit", func(c *gin.Context) {
		P.RunController(c, &http2.QuitGameController{})
	})

	router.POST("/api/game/role/assign", func(c *gin.Context) {
		P.RunController(c, &http2.AssignRoleController{})
	})

	router.POST("/api/game/name/assign", func(c *gin.Context) {
		P.RunController(c, &http2.AssignNameController{})
	})

	router.POST("/api/game/race/assign", func(c *gin.Context) {
		P.RunController(c, &http2.AssignRaceController{})
	})

	router.POST("/api/game/experience/adjust", func(c *gin.Context) {
		P.RunController(c, &http2.ChangePlayerExperienceController{})
	})

	router.POST("/api/game/health/adjust", func(c *gin.Context) {
		P.RunController(c, &http2.ChangePlayerHealthController{})
	})

	router.POST("/api/game/player/change/:index", func(c *gin.Context) {
		P.RunController(c, &http2.ChangePlayerNameController{})
	})

	router.POST("/api/game/equipment/add/:playerIndex", func(c *gin.Context) {
		P.RunController(c, &http2.AddCharacterEquipmentController{})
	})

	router.POST("/api/game/equipment/change/:playerIndex/:equipmentIndex", func(c *gin.Context) {
		P.RunController(c, &http2.ChangeCharacterEquipmentController{})
	})

	router.POST("/api/game/equipment/remove/:playerIndex/:equipmentIndex", func(c *gin.Context) {
		P.RunController(c, &http2.RemoveCharacterEquipmentController{})
	})

	router.POST("/api/game/name/change/:index", func(c *gin.Context) {
		P.RunController(c, &http2.ChangeCharacterNameController{})
	})

	router.POST("/api/game/race/change/:index", func(c *gin.Context) {
		P.RunController(c, &http2.ChangeCharacterRaceController{})
	})

	router.POST("/api/game/quest/change", func(c *gin.Context) {
		P.RunController(c, &http2.ChangeQuestController{})
	})

	router.GET("/api/game/options", func(c *gin.Context) {
		P.RunController(c, &http2.GetOptionsController{})
	})

	router.POST("/api/game/option/add", func(c *gin.Context) {
		P.RunController(c, &http2.AddOptionController{})
	})

	router.POST("/api/game/option/vote", func(c *gin.Context) {
		P.RunController(c, &http2.VoteOptionController{})
	})

	router.POST("/api/game/option/clear", func(c *gin.Context) {
		P.RunController(c, &http2.ClearRequestOptionsController{})
	})

	router.POST("/api/game/option/change/:index", func(c *gin.Context) {
		P.RunController(c, &http2.ChangeRequestOptionController{})
	})

	router.POST("/api/game/option/assign", func(c *gin.Context) {
		P.RunController(c, &http2.AssignResultController{})
	})

	router.PUT("/api/game/suggestion/new", func(c *gin.Context) {
		P.RunController(c, &http2.NewSuggestionController{})
	})

	router.POST("/api/game/dice/roll", func(c *gin.Context) {
		P.RunController(c, &http2.RollDiceController{})
	})

	router.POST("/api/game/suggestions", func(c *gin.Context) {
		P.RunController(c, &http2.ToggleSuggestionsController{})
	})

	router.POST("/api/game/volume", func(c *gin.Context) {
		P.RunController(c, &http2.SyncVolumeController{})
	})
}
