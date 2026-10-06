package middleware

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

// WebSocketUpgrade rejects anything that is not a websocket upgrade from an allowed origin,
// before the connection is upgraded. The websocket package itself allows every origin when
// none are configured, so this check is always required on websocket routes.
func WebSocketUpgrade() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		if !setup.OriginAllowed(c) {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
