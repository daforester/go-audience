package bindings

import (
	processors2 "bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"github.com/daforester/go-di-container/di"
)

type flag int64

const (
	FLAG_NONE flag = 0
)

func Register(app *di.App, flags flag) {
	app.Singleton((*repositories.GameRepository)(nil), app.Make(&repositories.MemoryGameRepository{}))
	app.Bind((*processors2.AssignRoleProcessorInterface)(nil), &processors2.AssignRoleProcessor{})
	app.Bind((*processors2.CreateGameProcessorInterface)(nil), &processors2.CreateGameProcessor{})
	app.Bind((*processors2.GetGameProcessorInterface)(nil), &processors2.GetGameProcessor{})
	app.Bind((*processors2.JoinGameProcessorInterface)(nil), &processors2.JoinGameProcessor{})
	app.Bind((*processors2.KickPlayerProcessorInterface)(nil), &processors2.KickPlayerProcessor{})
	app.Bind((*processors2.LoginProcessorInterface)(nil), &processors2.LoginProcessor{})
	app.Bind((*processors2.QuitGameProcessorInterface)(nil), &processors2.QuitGameProcessor{})
	app.Bind((*processors2.SyncVolumeProcessorInterface)(nil), &processors2.SyncVolumeProcessor{})
}
