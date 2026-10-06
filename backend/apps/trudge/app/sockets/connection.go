package sockets

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	uuid string
	websocket *websocket.Conn
	Related map[string]interface{}
	Tags []string
	mux *sync.Mutex
}

func (C Connection) New(socket *websocket.Conn) *Connection {
	c := new(Connection)
	c.websocket = socket
	c.Related = make(map[string]interface{})
	c.mux = new(sync.Mutex)
	return c
}

func (C *Connection) GetUUID() string {
	return C.uuid
}

func (C *Connection) SetUUID(uuid string) {
	C.uuid = uuid
}

func (C *Connection) SwapTags(from []string, to[]string) ([]string, error) {
	var err error
	if len(from) == len(to) {
		// 1 to 1 replacement
		for i, k := range from {
			for j, t := range C.Tags {
				if t == k {
					C.Tags[j] = to[i]
				}
			}
		}
	} else if len(to) == 1 {
		// Remove all From, Add Tag
		for _, k := range from {
			for j, t := range C.Tags {
				if t == k {
					C.Tags[j] = C.Tags[len(C.Tags)-1]
					C.Tags = C.Tags[:len(C.Tags)-1]
					break
				}
			}
		}
		C.Tags = append(C.Tags, to[0])
	} else if len(from) == 0 {
		// Add Tags
		C.Tags = append(C.Tags, to...)
	} else {
		err = errors.New("invalid replacement, from and to should be equal length or only one to length")
	}

	C.Tags = C.uniqueStrings(C.Tags)

	return C.Tags, err
}

func (C *Connection) WriteMessage(messageType int, data []byte) error {
	C.mux.Lock()
	defer C.mux.Unlock()
	return C.websocket.WriteMessage(messageType, data)
}

func (C *Connection) uniqueStrings(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}