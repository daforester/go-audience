package sockets

import (
	"github.com/gofiber/contrib/v3/websocket"
	"strings"
	"time"
)

// Clients PING every 30 seconds, a connection silent for longer than this is dead
const readWait = 90 * time.Second

type readFunc = func(*Collection, *Connection, string)

func ReadSocket(collection *Collection, connection *Connection, f ...readFunc) {
	for {
		_ = connection.websocket.SetReadDeadline(time.Now().Add(readWait))
		t, d, err := connection.websocket.ReadMessage()
		if err != nil {
			break
		}

		if t == websocket.TextMessage {
			text := string(d)
			if strings.HasPrefix(text, "PING") {
				_ = connection.WriteMessage(websocket.TextMessage, []byte("PONG"+text[4:]))
			}

			for n := 0; n < len(f); n++ {
				if f[n] != nil {
					f[n](collection, connection, text)
				}
			}
		}
	}
}
