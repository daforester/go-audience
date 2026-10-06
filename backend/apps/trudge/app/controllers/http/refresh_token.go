package http

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/repositories"
	"bitbucket.org/daforester/go-audience/services/controller"
	"bitbucket.org/daforester/go-audience/services/ginhelper"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
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
	c := C.GetContext().(*gin.Context)
	e := C.RefreshToken(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *RefreshTokenController) RefreshToken(c *gin.Context) *controller.Response {
	var tokens *processors.Tokens

	token, err := ginhelper.GetBearerToken(c, (&processors.LoginProcessor{}).ValidLoginToken)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if userID, ok := token.Claims.(jwt.MapClaims)["UID"]; ok {
		if userType, ok := token.Claims.(jwt.MapClaims)["UserType"]; ok {
			if "models.Player" == userType.(string) {
				ec, ok := token.Claims.(jwt.MapClaims)["Extra"]
				if !ok {
					return controller.ErrorResponse(http.StatusUnauthorized, "Extra claims not found")
				}
				extraClaims := ec.(map[string]interface{})
				rk, exists := extraClaims["RoomKey"]
				if !exists {
					return controller.ErrorResponse(http.StatusUnauthorized, "Room key not found")
				}
				roomKey := rk.(string)
				game := C.gameRepository.GetByKey(roomKey)
				if game == nil {
					return controller.ErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Game not found: %s", roomKey))
				}
				for _, player := range game.Players {
					if player.UUID == userID {
						tokens, err = C.loginProcessor.LoginTokens(player, userID.(string), true, extraClaims)
						break
					}
				}
				if tokens == nil {
					return controller.ErrorResponse(http.StatusUnauthorized, "User not found")
				}
			} else {
				return controller.ErrorResponse(http.StatusUnauthorized, "Unknown User Type")
			}
		} else {
			return controller.ErrorResponse(http.StatusUnauthorized, "Missing Token Type")
		}
	} else {
		return controller.ErrorResponse(http.StatusUnauthorized, "Missing UID")
	}

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, "Couldn't generate tokens")
	}

	return controller.StandardResponse(tokens)
}
