package sockets

import (
	"bitbucket.org/daforester/go-backend-repository/component/libs/strgen"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math"
	"runtime"
	"strconv"
	"sync"
)

type ConnectionClass int

const (
	CLIENT_ENDPOINT ConnectionClass = 0
	PEER_ENDPOINT   ConnectionClass = 1
)

type Collection struct {
	Name        string
	connections []*Connection
}

var WebSocketCollections = make(map[string]*Collection)
var lockCollections sync.Mutex

func NewCollection(name string) *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	ec, exists := WebSocketCollections[name]
	if exists {
		return ec
	}

	c := new(Collection)
	c.Name = name
	WebSocketCollections[c.Name] = c
	return c
}

func LockCollections() {
	lockCollections.Lock()
}

func UnlockCollections() {
	lockCollections.Unlock()
}

func RemoveCollection(name string) {
	LockCollections()
	defer UnlockCollections()

	_, exists := WebSocketCollections[name]
	if exists {
		delete(WebSocketCollections, name)
	}
}

func RemoveCollectionIfEmpty(name string) {
	c, exists := WebSocketCollections[name]
	if exists && len(c.connections) == 0 {
		RemoveCollection(name)
	}
}

func (C Collection) New() *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	c := new(Collection)
	c.Name, _ = strgen.RandString(5)
	_, exists := WebSocketCollections[c.Name]
	for exists {
		c.Name, _ = strgen.RandString(5)
		_, exists = WebSocketCollections[c.Name]
	}
	WebSocketCollections[c.Name] = c
	return c
}

func (C *Collection) AddConnection(connection *Connection) {
	C.connections = append(C.connections, connection)
}

func (C *Collection) RemoveConnection(connection *Connection) {
	for n, c := range C.connections {
		if c == connection {
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			break
		}
	}
}

func (C *Collection) RemoveConnectionByUUID(uuid string) {
	for n, c := range C.connections {
		if c.uuid == uuid {
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			break
		}
	}
}

func (C *Collection) UUIDExists(uuid string) bool {
	for _, c := range C.connections {
		if c.uuid == uuid {
			return true
		}
	}
	return false
}

func (C *Collection) ConnectionExists(connection *Connection) bool {
	for _, c := range C.connections {
		if c == connection {
			return true
		}
	}
	return false
}

func (C *Collection) Count() int {
	return len(C.connections)
}

func (C *Collection) String() string {
	output := strconv.Itoa(len(C.connections)) + " Connections\n"
	for _, conn := range C.connections {
		output += conn.uuid + ":" + conn.websocket.RemoteAddr().String() + ":" + conn.websocket.LocalAddr().String() + "\n"
	}
	return output
}

func (C *Collection) GetConnections() []*Connection {
	return C.connections
}

func (C *Collection) All() []*Connection {
	return C.connections
}

func (C *Collection) HasTag(tag string) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		for _, t := range c.Tags {
			if t == tag {
				a = append(a, c)
			}
		}
	}
	return a
}

func (C *Collection) HasRelated(key string, value interface{}) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		if c.Related[key] == value {
			a = append(a, c)
		}
	}
	return a
}

func (C *Collection) HasUUID(uuid string) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		if c.uuid == uuid {
			a = append(a, c)
		}
	}
	return a
}

func (C Collection) SendText(data string, connections ...*Connection) {
	C.Send(websocket.TextMessage, []byte(data), connections...)
}

func (C Collection) SendData(data []byte, connections ...*Connection) {
	C.Send(websocket.BinaryMessage, data, connections...)
}

func (C Collection) Send(messageType int, data []byte, connections ...*Connection) {
	if len(connections) == 0 {
		return
	}

	var wg sync.WaitGroup
	active := 0
	for _, c := range connections {
		if active >= int(math.Ceil(float64(runtime.NumCPU())/2)) {
			wg.Wait()
			active = 0
		}
		active += 1
		wg.Add(1)

		go func(connection *Connection) {
			defer wg.Done()
			err := connection.WriteMessage(messageType, data)
			if err != nil {
				logrus.Error(err.Error())
				_ = connection.websocket.Close()
			}
		}(c)
	}

	wg.Wait()
}
