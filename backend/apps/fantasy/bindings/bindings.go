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
	app.Bind((*processors2.ChangeSubmissionStatusProcessorInterface)(nil), &processors2.ChangeSubmissionStatusProcessor{})
	app.Bind((*processors2.CreateGameProcessorInterface)(nil), &processors2.CreateGameProcessor{})
	app.Bind((*processors2.GetGameProcessorInterface)(nil), &processors2.GetGameProcessor{})
	app.Bind((*processors2.JoinGameProcessorInterface)(nil), &processors2.JoinGameProcessor{})
	app.Bind((*processors2.KickPlayerProcessorInterface)(nil), &processors2.KickPlayerProcessor{})
	app.Bind((*processors2.LoginProcessorInterface)(nil), &processors2.LoginProcessor{})
	app.Bind((*processors2.NewSubmissionProcessorInterface)(nil), &processors2.NewSubmissionProcessor{})
	app.Bind((*processors2.QuitGameProcessorInterface)(nil), &processors2.QuitGameProcessor{})
	app.Bind((*processors2.ReactTopicProcessorInterface)(nil), &processors2.ReactTopicProcessor{})
	app.Bind((*processors2.SetGameCodeProcessorInterface)(nil), &processors2.SetGameCodeProcessor{})
	app.Bind((*processors2.SetPresentingTeamProcessorInterface)(nil), &processors2.SetPresentingTeamProcessor{})
	app.Bind((*processors2.StaffVoteProcessorInterface)(nil), &processors2.StaffVoteProcessor{})
	app.Bind((*processors2.StartRoundProcessorInterface)(nil), &processors2.StartRoundProcessor{})
	app.Bind((*processors2.SyncVolumeProcessorInterface)(nil), &processors2.SyncVolumeProcessor{})
	app.Bind((*processors2.TeamProcessorInterface)(nil), &processors2.TeamProcessor{})
	app.Bind((*processors2.ToggleSubmissionsProcessorInterface)(nil), &processors2.ToggleSubmissionsProcessor{})
	app.Bind((*processors2.ToggleSubmitterNamesProcessorInterface)(nil), &processors2.ToggleSubmitterNamesProcessor{})
	app.Bind((*processors2.VoteTeamProcessorInterface)(nil), &processors2.VoteTeamProcessor{})
}
