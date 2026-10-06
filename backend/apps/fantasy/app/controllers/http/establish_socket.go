package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/messages"
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/middleware"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"errors"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type EstablishSocketController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
}

func (C EstablishSocketController) New(gameRepository repositories.GameRepository) *EstablishSocketController {
	c := new(EstablishSocketController)
	c.gameRepository = gameRepository
	return c
}

func (C *EstablishSocketController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	return C.EstablishSocket(c)
}

// EstablishSocket validates the request before upgrading. Origin and upgrade headers have
// already been checked by middleware.WebSocketUpgrade on the route.
func (C *EstablishSocketController) EstablishSocket(c fiber.Ctx) error {
	roomKey, playerUUID, err := C.validateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(controller.ErrorResponse(fiber.StatusUnauthorized, err.Error()))
	}

	// The handler runs after this request has finished, so it only receives copies
	upgrade := websocket.New(func(conn *websocket.Conn) {
		C.serve(conn, roomKey, playerUUID)
	}, websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if r := recover(); r != nil {
				logrus.Error("websocket handler panic: ", r)
			}
		},
	})

	return upgrade(c)
}

func (C *EstablishSocketController) serve(conn *websocket.Conn, roomKey string, playerUUID string) {
	socketCollection := sockets.NewCollection(roomKey)

	ws := sockets.Connection{}.New(conn)
	ws.Related["room"] = roomKey
	ws.SetUUID(playerUUID)

	// Add the connection before reading the player's class under the game lock, so a role
	// change made at the same moment is applied to this connection's tags too
	socketCollection.AddConnection(ws)

	connected := C.welcome(socketCollection, ws, roomKey, playerUUID)
	if connected {
		sockets.ReadSocket(socketCollection, ws)
	}

	socketCollection.RemoveConnection(ws)
	_ = ws.Close()

	sockets.RemoveCollectionIfEmpty(socketCollection.Name)
}

func (C *EstablishSocketController) welcome(socketCollection *sockets.Collection, ws *sockets.Connection, roomKey string, playerUUID string) bool {
	game := C.gameRepository.GetByKey(roomKey)
	if game == nil {
		return false
	}

	_ = C.gameRepository.Lock(game)
	defer C.gameRepository.Unlock(game)

	player := game.GetPlayerByUUID(playerUUID)
	if player == nil {
		return false
	}

	ws.AddTags(string(player.Class))

	socketCollection.SendText("CODE "+strconv.Itoa(game.GameCode), ws)

	err := messages.ConnectedMessages(game, player)
	if err != nil {
		logrus.Error(err.Error())
	}

	return true
}

func (C *EstablishSocketController) validateRequest(c fiber.Ctx) (string, string, error) {
	roomKey := strings.ToUpper(fiberhelper.Param(c, "room"))
	playerUUID := fiberhelper.Param(c, "player")
	token := fiberhelper.Query(c, "token")

	if token == "" {
		return "", "", errors.New("token required")
	}

	tokenPlayer, err := middleware.ProcessPlayerTokenString(token)
	if err != nil {
		return "", "", err
	}

	if tokenPlayer.UUID != playerUUID || tokenPlayer.RoomKey != roomKey {
		return "", "", errors.New("token does not match player")
	}

	game := C.gameRepository.GetByKey(roomKey)
	if game == nil {
		return "", "", errors.New("game not found")
	}

	_ = C.gameRepository.Lock(game)
	player := game.GetPlayerByUUID(playerUUID)
	_ = C.gameRepository.Unlock(game)

	if player == nil {
		return "", "", errors.New("player not found")
	}

	return roomKey, playerUUID, nil
}
