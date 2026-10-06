package messages

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

// KickPlayerMessage tells one player's connections to clean themselves up, then closes them
// so a removed player stops receiving room broadcasts
func KickPlayerMessage(player *models.Player) {
	wsCollection := sockets.GetCollection(player.RoomKey)
	socks := wsCollection.HasUUID(player.UUID)
	wsCollection.SendText("QUIT", socks...)
	for _, sock := range socks {
		_ = sock.Close()
	}
}

// GameOverMessage tells every connection in the room the game has ended, then closes them
func GameOverMessage(roomKey string) {
	wsCollection := sockets.GetCollection(roomKey)
	socks := wsCollection.All()
	wsCollection.SendText("QUIT", socks...)
	for _, sock := range socks {
		_ = sock.Close()
	}
}
