package fiberhelper

import (
	"bytes"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Fiber hands out strings and byte slices that point into fasthttp's request buffers,
// which are reused as soon as the handler returns. Every value read from a request goes
// through these accessors so that what we keep is always a copy, never a reference.

func Param(c fiber.Ctx, key string) string {
	return strings.Clone(c.Params(key))
}

func Query(c fiber.Ctx, key string) string {
	return strings.Clone(c.Query(key))
}

func Header(c fiber.Ctx, key string) string {
	return strings.Clone(c.Get(key))
}

func Cookie(c fiber.Ctx, key string) string {
	return strings.Clone(c.Cookies(key))
}

func FormValue(c fiber.Ctx, key string) string {
	return strings.Clone(c.FormValue(key))
}

func Body(c fiber.Ctx) []byte {
	return bytes.Clone(c.Body())
}

// Host returns a copy of the host (with port) the request arrived on. X-Forwarded-Host is
// only honoured when the request came from a trusted proxy (fiber.Config.TrustProxy).
func Host(c fiber.Ctx) string {
	return strings.Clone(c.Host())
}

// BindJSON decodes the request body into out. The configured decoder (encoding/json)
// allocates fresh strings, so nothing in out refers to the request buffer.
func BindJSON(c fiber.Ctx, out interface{}) error {
	return c.Bind().JSON(out)
}
