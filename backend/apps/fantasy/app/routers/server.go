package routers

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/middleware"
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"encoding/json"
	"github.com/daforester/go-di-container/di"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/spf13/viper"
	"time"
)

// NewServer builds the Fiber app with Fantasy's security settings and all routes
func NewServer(app di.AppInterface) *fiber.App {
	config := fiber.Config{
		AppName: "fantasy " + setup.Version,

		// Request data is copied explicitly (see services/fiberhelper). Immutable stays off
		// so tests run against fasthttp's real buffer reuse and catch any missed copy.
		Immutable: false,

		// encoding/json allocates new strings when decoding, never aliasing the request
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,

		BodyLimit:    64 * 1024,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if viper.GetBool("TRUST_PROXY") {
		config.TrustProxy = true
		config.TrustProxyConfig = fiber.TrustProxyConfig{
			Proxies: viper.GetStringSlice("TRUSTED_PROXIES"),
		}
		config.ProxyHeader = fiber.HeaderXForwardedFor
	}

	r := fiber.New(config)

	r.Use(recover.New())

	r.Use(middleware.NaughtyCors())
	if corsConfig, ok := setup.CorsConfig(); ok {
		r.Use(cors.New(corsConfig))
	}

	NewAPIRouter(app, r)
	NewWebSocketRouter(app, r)

	return r
}
