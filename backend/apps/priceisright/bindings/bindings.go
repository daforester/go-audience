package bindings

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"github.com/daforester/go-di-container/di"
)

type flag int64

const (
	FLAG_NONE flag = 0
)

func Register(app *di.App, _ flag) {
	app.Bind((*processors.AssignRoleProcessorInterface)(nil), &processors.AssignRoleProcessor{})
	app.Bind((*processors.ChangeGamePasswordProcessorInterface)(nil), &processors.ChangeGamePasswordProcessor{})
	app.Bind((*processors.ChangePlayerBadgeProcessorInterface)(nil), &processors.ChangePlayerBadgeProcessor{})
	app.Bind((*processors.ChangePlayerNameProcessorInterface)(nil), &processors.ChangePlayerNameProcessor{})
	app.Singleton((*repositories.GameRepository)(nil), app.Make(&repositories.MemoryGameRepository{}))
	app.Bind((*processors.JoinGameProcessorInterface)(nil), &processors.JoinGameProcessor{})
	app.Bind((*processors.KickPlayerProcessorInterface)(nil), &processors.KickPlayerProcessor{})
	app.Bind((*processors.LoginProcessorInterface)(nil), &processors.LoginProcessor{})
	app.Bind((*processors.QuitGameProcessorInterface)(nil), &processors.QuitGameProcessor{})
	app.Bind((*processors.ChangeGuessProcessorInterface)(nil), &processors.ChangeGuessProcessor{})
	app.Bind((*processors.NewGuessProcessorInterface)(nil), &processors.NewGuessProcessor{})
	app.Bind((*processors.SaveImageUploadsProcessorInterface)(nil), &processors.SaveImageUploadsProcessor{})
	app.Bind((*processors.SaveRoundsProcessorInterface)(nil), &processors.SaveRoundsProcessor{})
	app.Bind((*processors.SaveImageUploadsProcessorInterface)(nil), &processors.SaveImageUploadsProcessor{})
	app.Bind((*processors.SetRoundProcessorInterface)(nil), &processors.SetRoundProcessor{})
	app.Bind((*processors.ClearGuessesProcessorInterface)(nil), &processors.ClearGuessesProcessor{})
	app.Bind((*processors.ActivateRoundProcessorInterface)(nil), &processors.ActivateRoundProcessor{})
	app.Bind((*processors.SummonPlayerProcessorInterface)(nil), &processors.SummonPlayerProcessor{})
}
