package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type AssignResultProcessorInterface interface {
	AssignResult(roomKey string, property string) (*models.Game, error)
}

type AssignResultProcessor struct {
	gameRepository repositories.GameRepository
}

func (P AssignResultProcessor) New(gameRepository repositories.GameRepository) *AssignResultProcessor {
	return &AssignResultProcessor{
		gameRepository,
	}
}

func (P *AssignResultProcessor) AssignResult(roomKey string, property string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	// Find top result
	best := uint(0)
	index := 0

	for i, vote := range game.Votes {
		if vote["Total"] >= best {
			index = i
			best = vote["Total"]
		}
	}

	if index >= len(game.Options) {
		return nil, errors.New("option not found")
	}

	option := game.Options[index]

	switch strings.Title(strings.ToLower(property)) {
	case "Weapon":
		player := game.GetPlayerByType(models.PLAYER_CLASS_WARRIOR)
		if player == nil {
			return game, errors.New("no warrior found")
		}
		if player.Equipment == nil {
			player.Equipment = make([]string, 0)
		}
		player.Equipment = append(player.Equipment, option)
		err := P.playerUpdateMessage(player)
		if err != nil {
			return game, err
		}
	case "Song":
		player := game.GetPlayerByType(models.PLAYER_CLASS_BARD)
		if player == nil {
			return game, errors.New("no bard found")
		}
		if player.Equipment == nil {
			player.Equipment = make([]string, 0)
		}
		player.Equipment = append(player.Equipment, option)
		err := P.playerUpdateMessage(player)
		if err != nil {
			return game, err
		}
	case "Spell":
		player := game.GetPlayerByType(models.PLAYER_CLASS_WIZARD)
		if player == nil {
			return game, errors.New("no wizard found")
		}
		if player.Equipment == nil {
			player.Equipment = make([]string, 0)
		}
		player.Equipment = append(player.Equipment, option)
		err := P.playerUpdateMessage(player)
		if err != nil {
			return game, err
		}
	case "Quest":
		game.Quest = option
		err := P.optionAssignedMessage(game.Key)
		if err != nil {
			return game, err
		}
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)
	game.Updated = time.Now()

	for playerIndex, player := range game.Players {
		if player.Vote == index {
			player.Experience += 1
			err := messages.ExperienceUpdateMessage(roomKey, player, uint(playerIndex), player.Experience)
			if err != nil {
				return game, err
			}
		}
		if player.Option == index {
			player.Experience += 3
			err := messages.ExperienceUpdateMessage(roomKey, player, uint(playerIndex), player.Experience)
			if err != nil {
				return game, err
			}
		}
	}

	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *AssignResultProcessor) optionAssignedMessage(roomKey string) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
	}{
		"OPTION_ASSIGNED",
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}

func (P *AssignResultProcessor) playerUpdateMessage(player *models.Player) error {
	wsCollection := sockets.NewCollection(player.RoomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    *models.Player
	}{
		"PLAYER_UPDATE",
		player,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasUUID(player.UUID)
	wsCollection.SendText(string(msg), socks...)

	return nil
}
