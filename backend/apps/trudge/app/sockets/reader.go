package sockets

import (
	"github.com/gorilla/websocket"
)

type readFunc = func(*Collection, *Connection, string)

func ReadSocket(collection *Collection, connection *Connection, f ...readFunc) {
	for {
		t, d, err := connection.websocket.ReadMessage()
		if err != nil {
			break
		}

		if t == websocket.TextMessage {
			text := string(d)
			if text[:4] == "PING" {
				_ = connection.websocket.WriteMessage(websocket.TextMessage, []byte("PONG" + text[4:]))
			}

			for n := 0; n < len(f); n ++ {
				if f[n] != nil {
					f[n](collection, connection, text)
				}
			}
		}
	}
}
