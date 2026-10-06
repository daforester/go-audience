package messages

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/sockets"
	"encoding/json"
	"errors"
)

func SuggestionsEnabledMessage(roomKey string, enabled bool) error {
	wsCollection := sockets.NewCollection(roomKey)

	msg, err := json.Marshal(struct {
		Command string
		Data    map[string]interface{}
	}{
		"SUGGESTIONS_ENABLED",
		map[string]interface{}{
			"Enabled": enabled,
		},
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	wsCollection.SendText(string(msg), wsCollection.All()...)

	return nil
}
