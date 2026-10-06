package middleware

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// NaughtyCors allows any origin whose host matches the host the request arrived on,
// when CORS_BE_NAUGHTY is enabled. Fiber's own CORS origin callback only receives the
// origin, so it cannot make this comparison. Never answers with a wildcard origin.
func NaughtyCors() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !viper.GetBool("CORS_BE_NAUGHTY") {
			return c.Next()
		}

		origin := fiberhelper.Header(c, fiber.HeaderOrigin)
		if origin == "" || !setup.OriginMatchesHost(origin, fiberhelper.Host(c)) {
			return c.Next()
		}

		c.Locals("NaughtyCors", true)
		c.Vary(fiber.HeaderOrigin)
		c.Set(fiber.HeaderAccessControlAllowOrigin, origin)
		c.Set(fiber.HeaderAccessControlAllowCredentials, "true")

		preflight := c.Method() == fiber.MethodOptions &&
			fiberhelper.Header(c, fiber.HeaderAccessControlRequestMethod) != ""

		if !preflight {
			c.Set(fiber.HeaderAccessControlExposeHeaders, strings.Join(setup.CorsExposeHeaders, ", "))
			return c.Next()
		}

		c.Vary(fiber.HeaderAccessControlRequestMethod, fiber.HeaderAccessControlRequestHeaders)
		c.Set(fiber.HeaderAccessControlAllowMethods, strings.Join(setup.CorsAllowMethods, ", "))
		c.Set(fiber.HeaderAccessControlAllowHeaders, strings.Join(setup.CorsAllowHeaders, ", "))
		c.Set(fiber.HeaderAccessControlMaxAge, strconv.Itoa(setup.CorsMaxAge))

		return c.SendStatus(fiber.StatusNoContent)
	}
}
