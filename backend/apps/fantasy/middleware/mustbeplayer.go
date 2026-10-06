package middleware

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func MustBePlayer() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := ProcessPlayerToken(c)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}
}

func ProcessPlayerToken(c fiber.Ctx) error {
	playerToken, err := fiberhelper.GetBearerToken(c, (&processors.LoginProcessor{}).ValidLoginToken)

	if err != nil {
		return err
	}

	player, err := PlayerFromToken(playerToken)
	if err != nil {
		return err
	}

	c.Locals("ClientToken", playerToken)
	c.Locals("Player", player)
	if extra, ok := playerToken.Claims.(jwt.MapClaims)["Extra"]; ok {
		c.Locals("Extra", extra)
	}

	return nil
}

// ProcessPlayerTokenString validates a raw access token, e.g. one passed on a websocket URL
func ProcessPlayerTokenString(token string) (*models.Player, error) {
	playerToken, err := (&processors.LoginProcessor{}).ValidLoginToken(token, true)
	if err != nil {
		return nil, err
	}

	return PlayerFromToken(playerToken)
}

func PlayerFromToken(playerToken *jwt.Token) (*models.Player, error) {
	claims, ok := playerToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	uType, ok := claims["UserType"].(string)
	if !ok || uType != "models.Player" {
		return nil, errors.New("unknown player type")
	}

	playerRaw, ok := claims["User"].(string)
	if !ok {
		return nil, errors.New("player not in token")
	}

	playerJson, err := base64.StdEncoding.DecodeString(playerRaw)
	if err != nil {
		return nil, err
	}

	player := new(models.Player)
	err = json.Unmarshal(playerJson, player)
	if err != nil {
		return nil, err
	}

	return player, nil
}
