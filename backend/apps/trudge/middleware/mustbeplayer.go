package middleware

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/app/processors"
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"bitbucket.org/daforester/go-audience/services/ginhelper"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MustBePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ProcessPlayerToken(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
		}
	}
}

func ProcessPlayerToken(c *gin.Context) error {
	playerToken, err := ginhelper.GetBearerToken(c, (&processors.LoginProcessor{}).ValidLoginToken)

	if err != nil {
		return err
	} else if uType, ok := playerToken.Claims.(jwt.MapClaims)["UserType"]; ok && uType.(string) == "models.Player" {
		player := new(models.Player)
		playerRaw := playerToken.Claims.(jwt.MapClaims)["User"].(string)

		playerJson, err := base64.StdEncoding.DecodeString(playerRaw)
		if err != nil {
			return err
		}

		err = json.Unmarshal(playerJson, player)
		if err != nil {
			return err
		}

		c.Set("ClientToken", playerToken)
		c.Set("Player", player)
		if extra, ok := playerToken.Claims.(jwt.MapClaims)["Extra"]; ok {
			c.Set("Extra", extra)
		}
	} else {
		return errors.New("unknown player type")
	}

	return nil
}
