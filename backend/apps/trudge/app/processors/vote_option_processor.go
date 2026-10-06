package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/messages"
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type VoteOptionProcessorInterface interface {
	VoteOption(roomKey string, player *models.Player, index int, option string) ([]map[string]uint, error)
}

type VoteOptionProcessor struct {
	gameRepository repositories.GameRepository
}

func (P VoteOptionProcessor) New(gameRepository repositories.GameRepository) *VoteOptionProcessor {
	return &VoteOptionProcessor{
		gameRepository,
	}
}

func (P *VoteOptionProcessor) VoteOption(roomKey string, player *models.Player, index int, option string) ([]map[string]uint, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	if len(game.Options) > len(game.Votes) {
		diff := make([]map[string]uint, len(game.Options)-len(game.Votes))
		game.Votes = append(game.Votes, diff...)
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	for i, o := range game.Options {
		if i == index && o == option {
			game.Votes[i]["Total"] += 1
			game.Votes[i][string(player.Class)] += 1

			for playerIndex, p := range game.Players {
				if p.UUID == player.UUID {
					p.Vote = i
				}
				if p.Option == i {
					p.Experience += 1
					err := messages.ExperienceUpdateMessage(roomKey, p, uint(playerIndex), p.Experience)
					if err != nil {
						return game.Votes, err
					}
				}
			}

			game.Updated = time.Now()
			err := P.gameRepository.Save(game)

			if err != nil {
				return game.Votes, err
			}

			err = P.optionVoteMessage(game.Key, player, index, option)
			if err != nil {
				return game.Votes, err
			}

			return game.Votes, nil
		}
	}

	return game.Votes, errors.New("voted option not found")
}

func (P *VoteOptionProcessor) optionVoteMessage(roomKey string, player *models.Player, index int, option string) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"OPTION_VOTE",
		map[string]interface{}{
			"index":  index,
			"option": option,
			"player": player,
		},
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
