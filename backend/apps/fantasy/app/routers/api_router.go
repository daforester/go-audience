package routers

import (
	http2 "bitbucket.org/daforester/go-audience/apps/fantasy/app/controllers/http"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/router"
	"github.com/daforester/go-di-container/di"
	"github.com/gofiber/fiber/v3"
)

type APIRouter struct {
	router.BaseRouter
}

func NewAPIRouter(app di.AppInterface, engine ...*fiber.App) router.Router {
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

func (P *APIRouter) RegisterRoutes(router *fiber.App) {
	P.RegisterPublicRoutes(router)
}

func (P *APIRouter) RegisterPublicRoutes(router *fiber.App) {
	router.Get("/api/token/refresh", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.RefreshTokenController{})
	})

	router.Get("/api/game", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.GetGameController{})
	})

	router.Post("/api/game/create", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.CreateGameController{})
	})

	router.Get("/api/game/kick/:uuid", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.KickPlayerController{})
	})

	router.Post("/api/game/join", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.JoinGameController{})
	})

	router.Get("/api/game/rejoin", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.RejoinGameController{})
	})

	router.Get("/api/game/quit", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.QuitGameController{})
	})

	router.Post("/api/game/role/assign", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.AssignRoleController{})
	})

	router.Post("/api/game/volume", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.SyncVolumeController{})
	})

	router.Post("/api/game/code", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.SetGameCodeController{})
	})

	router.Post("/api/game/team/add", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.AddTeamController{})
	})

	router.Post("/api/game/team/change/:id", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.ChangeTeamController{})
	})

	router.Post("/api/game/team/remove/:id", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.RemoveTeamController{})
	})

	router.Post("/api/game/submissions", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.ToggleSubmissionsController{})
	})

	router.Post("/api/game/submitter-names", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.ToggleSubmitterNamesController{})
	})

	router.Post("/api/game/round/start", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.StartRoundController{})
	})

	router.Post("/api/game/round/presenting", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.SetPresentingTeamController{})
	})

	router.Put("/api/game/submission/new", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.NewSubmissionController{})
	})

	router.Post("/api/game/submission/status/:id", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.ChangeSubmissionStatusController{})
	})

	router.Post("/api/game/submission/vote/:id", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.StaffVoteController{})
	})

	router.Post("/api/game/round/react", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.ReactTopicController{})
	})

	router.Post("/api/game/round/vote", func(c fiber.Ctx) error {
		return P.RunController(c, &http2.VoteTeamController{})
	})
}
