package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
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

func (P *QuitGameProcessor) QuitGame(player *models.Player) error {
	game := P.gameRepository.GetByKey(player.RoomKey)
	if game == nil {
		return errors.New("no game with room key")
	}

	otherLeadersCount := 0
	for _, p := range game.Players {
		if p.Class == models.PLAYER_CLASS_LEADER && p.UUID != player.UUID {
			otherLeadersCount += 1
		}
	}

	err := P.playerQuitMessage(player)

	if err != nil {
		return err
	}

	if player.Class == models.PLAYER_CLASS_LEADER && otherLeadersCount == 0 {
		// Destroy Game
		P.gameOverMessage(game)
		_ = P.gameRepository.Lock(nil)
		defer P.gameRepository.Unlock(nil)
		P.gameRepository.Delete(game)
		nowTime := time.Now()
		game.Deleted = &nowTime
	} else {
		// Remove Player
		_ = P.gameRepository.Lock(game)
		defer P.gameRepository.Unlock(game)
		game.Updated = time.Now()
		game.RemovePlayer(player)
		err = P.gameRepository.Save(game)
	}

	if err != nil {
		return err
	}

	return nil
}

func (P *QuitGameProcessor) playerQuitMessage(player *models.Player) error {
	msg, err := json.Marshal(struct {
		Command string
		Data    *models.Player
	}{
		"PLAYER_QUIT",
		player,
	})
	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	wsCollection := sockets.NewCollection(player.RoomKey)

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)
	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}

func (P *QuitGameProcessor) gameOverMessage(game *models.Game) {
	wsCollection := sockets.NewCollection(game.Key)
	wsCollection.SendText("QUIT", wsCollection.All()...)
	// Clients do not need to send anything to server after receiving this, just clean themselves up
}
