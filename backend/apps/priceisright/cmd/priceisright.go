package main

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/routers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/bindings"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/apps/priceisright/setup"
	"bitbucket.org/daforester/go-audience/services/sessions"
	"github.com/daforester/go-di-container/di"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

func main() {
	viper.SetDefault("ENVIRONMENT", "dev")

	setup.Config()
	setup.Keys()
	setup.Logger()
	setup.MimeTypes()

	app := di.New()

	bindings.Register(app, bindings.FLAG_NONE)

	repo := app.Make((*repositories.GameRepository)(nil))
	gameMem := repo.(*repositories.MemoryGameRepository)
	_ = gameMem.LoadFromDisk(viper.GetString("GAME_STATE_FILE"))

	game := gameMem.Get()
	rounds := make([]*models.Round, 0)
	for _, round := range game.Rounds {
		if round != nil {
			rounds = append(rounds, round)
		}
	}
	game.Rounds = rounds
	gameMem.Save(game)

	r := gin.Default()

	err := sessions.SetupSessionUsingViper(r, "sessions")
	if err != nil {
		panic(err)
	}

	r.Use(location.New(location.Config{
		Scheme: viper.GetString("HTTP_SCHEMA"),
		Host:   viper.GetString("HTTP_HOST"),
		Base:   viper.GetString("HTTP_BASE"),
	}))

	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:5173", "https://priceisright.show", "https://api.priceisright.show", "http://localhost.priceisright.show"},
		AllowMethods:     []string{"DELETE", "GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	routers.NewViewRouter(app, r)
	routers.NewAPIRouter(app, r)
	routers.NewWebSocketRouter(app, r)

	go func() {
		_ = r.Run(":" + viper.GetString("HTTP_PORT")) // listen and serve on 0.0.0.0:8080
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	_ = gameMem.SaveToDisk(viper.GetString("GAME_STATE_FILE"))
}
