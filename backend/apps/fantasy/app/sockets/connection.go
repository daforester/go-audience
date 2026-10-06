package sockets

import (
	"errors"
	"github.com/gofiber/contrib/v3/websocket"
	"sync"
	"time"
)

const writeWait = 5 * time.Second

// MaxMessageSize is the largest message a client may send. Clients only send PINGs.
const MaxMessageSize = 4096

type Connection struct {
	uuid      string
	websocket *websocket.Conn
	Related   map[string]interface{}
	Tags      []string
	mux       *sync.Mutex
	tagMux    *sync.RWMutex
}

func (C Connection) New(socket *websocket.Conn) *Connection {
	c := new(Connection)
	c.websocket = socket
	c.Related = make(map[string]interface{})
	c.mux = new(sync.Mutex)
	c.tagMux = new(sync.RWMutex)
	if socket != nil {
		socket.SetReadLimit(MaxMessageSize)
	}
	return c
}

func (C *Connection) GetUUID() string {
	return C.uuid
}

func (C *Connection) SetUUID(uuid string) {
	C.uuid = uuid
}

func (C *Connection) AddTags(tags ...string) []string {
	C.tagMux.Lock()
	defer C.tagMux.Unlock()
	C.Tags = C.uniqueStrings(append(C.Tags, tags...))
	return C.copyTags()
}

func (C *Connection) HasTag(tag string) bool {
	C.tagMux.RLock()
	defer C.tagMux.RUnlock()
	for _, t := range C.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (C *Connection) SwapTags(from []string, to []string) ([]string, error) {
	C.tagMux.Lock()
	defer C.tagMux.Unlock()

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

	return C.copyTags(), err
}

func (C *Connection) WriteMessage(messageType int, data []byte) error {
	C.mux.Lock()
	defer C.mux.Unlock()
	_ = C.websocket.SetWriteDeadline(time.Now().Add(writeWait))
	return C.websocket.WriteMessage(messageType, data)
}

func (C *Connection) Close() error {
	return C.websocket.Close()
}

func (C *Connection) copyTags() []string {
	t := make([]string, len(C.Tags))
	copy(t, C.Tags)
	return t
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
