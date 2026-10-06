package sockets

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/libs/strgen"
	"github.com/gofiber/contrib/v3/websocket"
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
	mux         sync.RWMutex
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

// GetCollection returns the room's collection without registering one. Use it to send:
// a room nobody is connected to gets an empty collection that is never stored.
func GetCollection(name string) *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	ec, exists := WebSocketCollections[name]
	if exists {
		return ec
	}

	c := new(Collection)
	c.Name = name
	return c
}

// JoinCollection adds the connection to the room's collection, creating it if needed.
// Done under the collections lock so it cannot race LeaveCollection removing the room.
func JoinCollection(name string, connection *Connection) *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	c, exists := WebSocketCollections[name]
	if !exists {
		c = new(Collection)
		c.Name = name
		WebSocketCollections[name] = c
	}
	c.AddConnection(connection)
	return c
}

// LeaveCollection removes the connection, and the room's collection once it is empty
func LeaveCollection(collection *Collection, connection *Connection) {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	collection.RemoveConnection(connection)
	if collection.Count() == 0 && WebSocketCollections[collection.Name] == collection {
		delete(WebSocketCollections, collection.Name)
	}
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
	LockCollections()
	defer UnlockCollections()

	c, exists := WebSocketCollections[name]
	if exists && c.Count() == 0 {
		delete(WebSocketCollections, name)
	}
}

func (C *Collection) New() *Collection {
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
	C.mux.Lock()
	defer C.mux.Unlock()
	C.connections = append(C.connections, connection)
}

func (C *Collection) RemoveConnection(connection *Connection) {
	C.mux.Lock()
	defer C.mux.Unlock()
	for n, c := range C.connections {
		if c == connection {
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			break
		}
	}
}

func (C *Collection) RemoveConnectionByUUID(uuid string) {
	C.mux.Lock()
	defer C.mux.Unlock()
	for n, c := range C.connections {
		if c.uuid == uuid {
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			break
		}
	}
}

func (C *Collection) UUIDExists(uuid string) bool {
	return len(C.HasUUID(uuid)) > 0
}

func (C *Collection) ConnectionExists(connection *Connection) bool {
	C.mux.RLock()
	defer C.mux.RUnlock()
	for _, c := range C.connections {
		if c == connection {
			return true
		}
	}
	return false
}

func (C *Collection) Count() int {
	C.mux.RLock()
	defer C.mux.RUnlock()
	return len(C.connections)
}

func (C *Collection) String() string {
	connections := C.All()
	output := strconv.Itoa(len(connections)) + " Connections\n"
	for _, conn := range connections {
		output += conn.uuid + ":" + conn.websocket.RemoteAddr().String() + ":" + conn.websocket.LocalAddr().String() + "\n"
	}
	return output
}

func (C *Collection) GetConnections() []*Connection {
	return C.All()
}

// All returns a copy of the connection list, safe to range over while others connect
func (C *Collection) All() []*Connection {
	C.mux.RLock()
	defer C.mux.RUnlock()
	a := make([]*Connection, len(C.connections))
	copy(a, C.connections)
	return a
}

func (C *Collection) HasTag(tag string) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.All() {
		if c.HasTag(tag) {
			a = append(a, c)
		}
	}
	return a
}

// HasTags returns connections carrying any of the tags, each connection at most once
func (C *Collection) HasTags(tags ...string) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.All() {
		for _, t := range tags {
			if c.HasTag(t) {
				a = append(a, c)
				break
			}
		}
	}
	return a
}

func (C *Collection) HasRelated(key string, value interface{}) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.All() {
		if c.Related[key] == value {
			a = append(a, c)
		}
	}
	return a
}

func (C *Collection) HasUUID(uuid string) []*Connection {
	a := make([]*Connection, 0)
	for _, c := range C.All() {
		if c.uuid == uuid {
			a = append(a, c)
		}
	}
	return a
}

func (C *Collection) SendText(data string, connections ...*Connection) {
	C.Send(websocket.TextMessage, []byte(data), connections...)
}

func (C *Collection) SendData(data []byte, connections ...*Connection) {
	C.Send(websocket.BinaryMessage, data, connections...)
}

func (C *Collection) Send(messageType int, data []byte, connections ...*Connection) {
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
