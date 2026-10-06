package setup

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

var CorsAllowMethods = []string{fiber.MethodDelete, fiber.MethodGet, fiber.MethodPost, fiber.MethodPut}
var CorsAllowHeaders = []string{fiber.HeaderOrigin, fiber.HeaderContentType, "X-Auth-Token", fiber.HeaderAuthorization}
var CorsExposeHeaders = []string{fiber.HeaderContentLength}

const CorsMaxAge = 12 * 60 * 60

func Cors() {
	if viper.GetBool("CORS_BE_NAUGHTY") {
		logrus.Warn("CORS_BE_NAUGHTY is on: any origin matching the request host is allowed")
	} else if len(AllowedOrigins()) == 0 {
		logrus.Warn("CORS_ALLOW_ORIGINS is empty: browsers will only reach the API from its own origin")
	}
}

// AllowedOrigins are the explicitly configured origins
func AllowedOrigins() []string {
	origins := make([]string, 0)
	for _, o := range viper.GetStringSlice("CORS_ALLOW_ORIGINS") {
		o = strings.TrimRight(strings.TrimSpace(o), "/")
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

// OriginAllowed is the single origin check shared by HTTP CORS and the websocket upgrade
func OriginAllowed(c fiber.Ctx) bool {
	origin := fiberhelper.Header(c, fiber.HeaderOrigin)
	if origin == "" {
		return false
	}

	for _, o := range AllowedOrigins() {
		if strings.EqualFold(o, origin) {
			return true
		}
	}

	if viper.GetBool("CORS_BE_NAUGHTY") {
		return OriginMatchesHost(origin, fiberhelper.Host(c))
	}

	return false
}

// OriginMatchesHost reports whether the origin points at the host the request arrived on
func OriginMatchesHost(origin string, host string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" || host == "" {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return strings.EqualFold(u.Host, host)
}

// CorsConfig is the allow-list CORS middleware config. Requests already answered by the
// naughty CORS middleware are skipped.
func CorsConfig() (cors.Config, bool) {
	origins := AllowedOrigins()
	return cors.Config{
		Next: func(c fiber.Ctx) bool {
			handled, _ := c.Locals("NaughtyCors").(bool)
			return handled
		},
		AllowCredentials: true,
		AllowOrigins:     origins,
		AllowMethods:     CorsAllowMethods,
		AllowHeaders:     CorsAllowHeaders,
		ExposeHeaders:    CorsExposeHeaders,
		MaxAge:           CorsMaxAge,
	}, len(origins) > 0
}
