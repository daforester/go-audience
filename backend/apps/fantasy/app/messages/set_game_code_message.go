package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"strconv"
)

func SetGameCodeMessage(roomKey string, code int) {
	wsCollection := sockets.GetCollection(roomKey)
	wsCollection.SendText("CODE "+strconv.Itoa(code), wsCollection.All()...)
}
