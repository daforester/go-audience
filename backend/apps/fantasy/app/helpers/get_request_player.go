package helpers

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/middleware"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"errors"
	"github.com/gofiber/fiber/v3"
)

// GetRequestPlayer returns the player from the request's access token. It only identifies
// the requester (UUID, RoomKey); processors look up the stored player for anything else,
// such as role checks, since the role may have changed since the token was issued.
func GetRequestPlayer(c fiber.Ctx) (*models.Player, error) {
	player, ok := c.Locals("Player").(*models.Player)
	if !ok {
		err := middleware.ProcessPlayerToken(c)
		if err != nil {
			return nil, errors.New("player token data not provided")
		}
		player, ok = c.Locals("Player").(*models.Player)
		if !ok {
			return nil, errors.New("player token data not parsed")
		}
	}

	return player, nil
}
