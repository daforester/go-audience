package main

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/routers"
	"bitbucket.org/daforester/go-audience/apps/fantasy/bindings"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"github.com/daforester/go-di-container/di"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	setup.Config()
	setup.Keys()
	setup.Logger()
	setup.Cors()

	app := di.New()

	bindings.Register(app, bindings.FLAG_NONE)

	repo := app.Make((*repositories.GameRepository)(nil))
	gameMem := repo.(*repositories.MemoryGameRepository)
	_ = gameMem.LoadFromDisk(viper.GetString("GAME_STATE_FILE"))

	r := routers.NewServer(app)

	go func() {
		err := r.Listen(":" + viper.GetString("HTTP_PORT"))
		if err != nil {
			logrus.Error(err.Error())
		}
	}()

	// Snapshot the games regularly so a crash loses at most one interval of play
	stopSnapshots := make(chan struct{})
	go func() {
		interval := time.Duration(viper.GetInt("SNAPSHOT_INTERVAL_SECONDS")) * time.Second
		if interval <= 0 {
			return
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := gameMem.SaveToDisk(viper.GetString("GAME_STATE_FILE")); err != nil {
					logrus.Error(err.Error())
				}
			case <-stopSnapshots:
				return
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	close(stopSnapshots)
	_ = r.ShutdownWithTimeout(5 * time.Second)

	_ = gameMem.SaveToDisk(viper.GetString("GAME_STATE_FILE"))
}
