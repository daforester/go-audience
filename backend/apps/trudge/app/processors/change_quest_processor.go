package processors

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"encoding/json"
	"errors"
	"time"
)

type ChangeQuestProcessorInterface interface {
	ChangeQuest(roomKey string, value string) (*models.Game, error)
}

type ChangeQuestProcessor struct {
	gameRepository repositories.GameRepository
}

func (P ChangeQuestProcessor) New(gameRepository repositories.GameRepository) *ChangeQuestProcessor {
	return &ChangeQuestProcessor{
		gameRepository,
	}
}

func (P *ChangeQuestProcessor) ChangeQuest(roomKey string, value string) (*models.Game, error) {
	game := P.gameRepository.GetByKey(roomKey)
	if game == nil {
		return nil, errors.New("game not found")
	}

	_ = P.gameRepository.Lock(game)
	defer P.gameRepository.Unlock(game)

	game.Quest = value
	game.Updated = time.Now()
	err := P.gameRepository.Save(game)

	if err != nil {
		return game, err
	}

	err = P.questChangedMessage(roomKey, value)
	if err != nil {
		return game, err
	}

	return game, nil
}

func (P *ChangeQuestProcessor) questChangedMessage(roomKey string, value string) error {
	wsCollection := sockets.NewCollection(roomKey)

	data := struct {
		Quest string
	}{
		value,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"QUEST_CHANGED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.All()
	wsCollection.SendText(string(msg), socks...)

	return nil
}
