package bindings

import (
	processors2 "bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"github.com/daforester/go-di-container/di"
)

type flag int64

const (
	FLAG_NONE flag = 0
)

func Register(app *di.App, flags flag) {
	// Has a trudge ID which could be changed per request, so must be unique instance
	app.Singleton((*repositories.GameRepository)(nil), app.Make(&repositories.MemoryGameRepository{}))
	app.Bind((*processors2.AddOptionProcessorInterface)(nil), &processors2.AddOptionProcessor{})
	app.Bind((*processors2.AddCharacterEquipmentProcessorInterface)(nil), &processors2.AddCharacterEquipmentProcessor{})
	app.Bind((*processors2.AssignNameProcessorInterface)(nil), &processors2.AssignNameProcessor{})
	app.Bind((*processors2.AssignRaceProcessorInterface)(nil), &processors2.AssignRaceProcessor{})
	app.Bind((*processors2.AssignResultProcessorInterface)(nil), &processors2.AssignResultProcessor{})
	app.Bind((*processors2.AssignRoleProcessorInterface)(nil), &processors2.AssignRoleProcessor{})
	app.Bind((*processors2.ChangeCharacterEquipmentProcessorInterface)(nil), &processors2.ChangeCharacterEquipmentProcessor{})
	app.Bind((*processors2.ChangeCharacterNameProcessorInterface)(nil), &processors2.ChangeCharacterNameProcessor{})
	app.Bind((*processors2.ChangeCharacterRaceProcessorInterface)(nil), &processors2.ChangeCharacterRaceProcessor{})
	app.Bind((*processors2.ChangePlayerNameProcessorInterface)(nil), &processors2.ChangePlayerNameProcessor{})
	app.Bind((*processors2.ChangePlayerExperienceProcessorInterface)(nil), &processors2.ChangePlayerExperienceProcessor{})
	app.Bind((*processors2.ChangePlayerHealthProcessorInterface)(nil), &processors2.ChangePlayerHealthProcessor{})
	app.Bind((*processors2.ChangeQuestProcessorInterface)(nil), &processors2.ChangeQuestProcessor{})
	app.Bind((*processors2.ChangeRequestOptionProcessorInterface)(nil), &processors2.ChangeRequestOptionProcessor{})
	app.Bind((*processors2.ClearRequestOptionsProcessorInterface)(nil), &processors2.ClearRequestOptionsProcessor{})
	app.Bind((*processors2.CreateGameProcessorInterface)(nil), &processors2.CreateGameProcessor{})
	app.Bind((*processors2.DiceRollProcessorInterface)(nil), &processors2.DiceRollProcessor{})
	app.Bind((*processors2.JoinGameProcessorInterface)(nil), &processors2.JoinGameProcessor{})
	app.Bind((*processors2.KickPlayerProcessorInterface)(nil), &processors2.KickPlayerProcessor{})
	app.Bind((*processors2.LoginProcessorInterface)(nil), &processors2.LoginProcessor{})
	app.Bind((*processors2.NewSuggestionProcessorInterface)(nil), &processors2.NewSuggestionProcessor{})
	app.Bind((*processors2.QuitGameProcessorInterface)(nil), &processors2.QuitGameProcessor{})
	app.Bind((*processors2.RemoveCharacterEquipmentProcessorInterface)(nil), &processors2.RemoveCharacterEquipmentProcessor{})
	app.Bind((*processors2.SetGameCodeProcessorInterface)(nil), &processors2.SetGameCodeProcessor{})
	app.Bind((*processors2.SyncVolumeProcessorInterface)(nil), &processors2.SyncVolumeProcessor{})
	app.Bind((*processors2.ToggleSuggestionsProcessorInterface)(nil), &processors2.ToggleSuggestionsProcessor{})
	app.Bind((*processors2.VoteOptionProcessorInterface)(nil), &processors2.VoteOptionProcessor{})
}
