package http

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/fiberhelper"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenController struct {
	controller.HTTP
	gameRepository repositories.GameRepository
	loginProcessor processors.LoginProcessorInterface
}

func (C RefreshTokenController) New(gameRepository repositories.GameRepository, loginProcessor processors.LoginProcessorInterface) *RefreshTokenController {
	c := new(RefreshTokenController)
	c.gameRepository = gameRepository
	c.loginProcessor = loginProcessor
	return c
}

func (C *RefreshTokenController) Run() error {
	c := C.GetContext().(fiber.Ctx)
	e := C.RefreshToken(c)

	return c.Status(e.ResponseCode).JSON(e)
}

func (C *RefreshTokenController) RefreshToken(c fiber.Ctx) *controller.Response {
	token, err := fiberhelper.GetBearerToken(c, C.loginProcessor.ValidLoginToken)

	if err != nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["UID"].(string)
	if !ok {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "Missing UID")
	}

	userType, ok := claims["UserType"].(string)
	if !ok {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "Missing Token Type")
	}

	if userType != "models.Player" {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "Unknown User Type")
	}

	extraClaims, ok := claims["Extra"].(map[string]interface{})
	if !ok {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "Extra claims not found")
	}

	roomKey, ok := extraClaims["RoomKey"].(string)
	if !ok {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "Room key not found")
	}

	game := C.gameRepository.GetByKey(roomKey)
	if game == nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, fmt.Sprintf("Game not found: %s", roomKey))
	}

	var player *models.Player
	_ = C.gameRepository.Lock(game)
	if p := game.GetPlayerByUUID(userID); p != nil {
		player = p.Copy()
	}
	_ = C.gameRepository.Unlock(game)

	if player == nil {
		return controller.ErrorResponse(fiber.StatusUnauthorized, "User not found")
	}

	tokens, err := C.loginProcessor.LoginTokens(player, userID, true, extraClaims)
	if err != nil {
		return controller.ErrorResponse(fiber.StatusInternalServerError, "Couldn't generate tokens")
	}

	return controller.StandardResponse(tokens)
}
