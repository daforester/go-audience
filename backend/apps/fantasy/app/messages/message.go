package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"encoding/json"
	"errors"
)

// sendToClasses marshals a command message and sends it to every connection tagged with
// one of the player classes
func sendToClasses(roomKey string, command string, data interface{}, classes ...models.PlayerClass) error {
	msg, err := marshal(command, data)
	if err != nil {
		return err
	}

	wsCollection := sockets.GetCollection(roomKey)

	tags := make([]string, 0, len(classes))
	for _, class := range classes {
		tags = append(tags, string(class))
	}

	socks := wsCollection.HasTags(tags...)
	wsCollection.SendText(msg, socks...)

	return nil
}

func sendToPlayer(roomKey string, uuid string, command string, data interface{}) error {
	msg, err := marshal(command, data)
	if err != nil {
		return err
	}

	wsCollection := sockets.GetCollection(roomKey)
	socks := wsCollection.HasUUID(uuid)
	wsCollection.SendText(msg, socks...)

	return nil
}

func sendToAll(roomKey string, command string, data interface{}) error {
	msg, err := marshal(command, data)
	if err != nil {
		return err
	}

	wsCollection := sockets.GetCollection(roomKey)
	wsCollection.SendText(msg, wsCollection.All()...)

	return nil
}

func marshal(command string, data interface{}) (string, error) {
	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		command,
		data,
	})

	if err != nil {
		return "", errors.New("unable to marshal socket message")
	}

	return string(msg), nil
}
