package ginhelper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetBearerToken(c *gin.Context, validatorFunc func(string, bool) (*jwt.Token, error)) (*jwt.Token, error) {
	header := c.Request.Header["Authorization"]

	if len(header) > 0 {
		bits := strings.Split(header[0], " ")
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
