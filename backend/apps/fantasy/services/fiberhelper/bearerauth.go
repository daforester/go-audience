package fiberhelper

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func GetBearerToken(c fiber.Ctx, validatorFunc func(string, bool) (*jwt.Token, error)) (*jwt.Token, error) {
	header := Header(c, fiber.HeaderAuthorization)

	if header != "" {
		bits := strings.Split(header, " ")
		if len(bits) == 2 {
			t := bits[0]
			d := bits[1]

			if strings.ToLower(t) == "bearer" {
				return validatorFunc(d, true)
			}
		}
	}

	return nil, errors.New("token not found")
}
