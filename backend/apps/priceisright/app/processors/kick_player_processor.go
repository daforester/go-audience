package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"time"
)

type KickPlayerProcessorInterface interface {
	KickPlayer(leader *models.Player, uuid string) error
}

type KickPlayerProcessor struct {
	gameRepository repositories.GameRepository
}

func (P KickPlayerProcessor) New(gameRepository repositories.GameRepository) *KickPlayerProcessor {
	return &KickPlayerProcessor{
		gameRepository,
	}
}

func (P *KickPlayerProcessor) KickPlayer(leader *models.Player, uuid string) error {
	if leader.Class != models.PLAYER_CLASS_LEADER {
		return errors.New("must be room leader to kick players")
	}

	game := P.gameRepository.Get()
	if game == nil {
		return errors.New("no game with room key")
	}

	var player *models.Player

	for _, p := range game.Players {
		if p.UUID == uuid {
			player = p
		}
	}

	if player == nil {
		return errors.New("player not found")
	}

	P.kickPlayerMessage(player)

	err := P.playerQuitMessage(player)

	if err != nil {
		return err
	}

	P.gameRepository.Lock()
	defer P.gameRepository.Unlock()
	game.Updated = time.Now()
	game.RemovePlayer(player)
	P.gameRepository.Save(game)

	return nil
}

func (P *KickPlayerProcessor) playerQuitMessage(player *models.Player) error {
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

	wsCollection := sockets.NewCollection("priceIsRight")

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)
	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}

func (P *KickPlayerProcessor) kickPlayerMessage(player *models.Player) {
	wsCollection := sockets.NewCollection("priceIsRight")
	socks := wsCollection.HasTag(string(player.Class))
	for _, sock := range socks {
		p, exists := sock.Related["player"]
		if exists {
			sockPlayer := p.(*models.Player)
			if sockPlayer.UUID == player.UUID {
				wsCollection.SendText("QUIT", sock)
				break
			}
		}
	}
	// Clients do not need to send anything to server after receiving this, just clean themselves up
}
