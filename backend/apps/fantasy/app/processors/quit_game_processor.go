package processors

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"time"
)

type QuitGameProcessorInterface interface {
	QuitGame(player *models.Player) error
}

type QuitGameProcessor struct {
	gameRepository repositories.GameRepository
}

func (P QuitGameProcessor) New(gameRepository repositories.GameRepository) *QuitGameProcessor {
	return &QuitGameProcessor{
		gameRepository,
	}
}

func (P *QuitGameProcessor) QuitGame(requester *models.Player) error {
	game := P.gameRepository.GetByKey(requester.RoomKey)
	if game == nil {
		return ErrGameNotFound
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(requester.UUID)
	if player == nil {
		return ErrPlayerNotFound
	}

	otherAdminsCount := 0
	for _, p := range game.GetAdmins() {
		if p.UUID != player.UUID {
			otherAdminsCount += 1
		}
	}

	err := messages.PlayerQuitMessage(player)

	if err != nil {
		return err
	}

	if player.Class == models.PLAYER_CLASS_ADMIN && otherAdminsCount == 0 {
		// Destroy Game
		messages.GameOverMessage(game.Key)
		err = P.gameRepository.Delete(game)
		nowTime := time.Now()
		game.Deleted = &nowTime
	} else {
		// Remove Player
		messages.KickPlayerMessage(player)
		game.Updated = time.Now()
		game.RemovePlayer(player)
		err = P.gameRepository.Save(game)
		if err == nil && player.Class == models.PLAYER_CLASS_DIRECTOR {
			err = messages.GameUpdateMessage(game)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
