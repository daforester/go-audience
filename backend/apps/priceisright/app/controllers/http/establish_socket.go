package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	c := C.GetContext().(*gin.Context)
	return C.EstablishSocket(c)
}

func (C *EstablishSocketController) EstablishSocket(c *gin.Context) error {
	var wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	w := c.Writer
	r := c.Request

	socket, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		logrus.Error("Failed to set websocket upgrade: %+v", err)
		return err
	}

	game, player, err := C.validateRequest(c)
	if err != nil {
		socket.Close()
		return nil
	}

	socketCollection := sockets.NewCollection("priceIsRight")

	ws := sockets.Connection{}.New(socket)
	ws.Related["game"] = game
	ws.Related["player"] = player
	ws.Tags = append(ws.Tags, string(player.Class))
	ws.SetUUID(player.UUID)

	socketCollection.AddConnection(ws)

	socketCollection.SendText("CODE "+strconv.Itoa(game.GameCode), ws)

	sockets.ReadSocket(socketCollection, ws)

	socketCollection.RemoveConnection(ws)

	sockets.RemoveCollectionIfEmpty(socketCollection.Name)

	return nil
}

func (C *EstablishSocketController) validateRequest(c *gin.Context) (*models.Game, *models.Player, error) {
	playerUUID := c.Param("player")

	game := C.gameRepository.Get()
	if game == nil {
		return nil, nil, errors.New("game not found")
	}

	var player *models.Player

	for _, p := range game.Players {
		if p.GetUUID() == playerUUID {
			player = p
		}
	}

	if player == nil {
		return nil, nil, errors.New("player not found")
	}

	return game, player, nil
}
