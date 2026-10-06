package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/middleware"
	"encoding/base32"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

type Tokens struct {
	AccessToken         string `json:"AccessToken"`
	AccessTokenExpires  int64  `json:"AccessExpires"`
	RefreshToken        string `json:"RefreshToken"`
	RefreshTokenExpires int64  `json:"RefreshExpires"`
}

type LoginProcessorInterface interface {
	LoginTokens(user interface{}, uid string, encryptToken bool, claims map[string]interface{}) (*Tokens, error)
}

type LoginProcessor struct {
}

func (P LoginProcessor) New() *LoginProcessor {
	return new(LoginProcessor)
}

func (P *LoginProcessor) LoginTokens(user interface{}, uid string, encryptToken bool, claims map[string]interface{}) (*Tokens, error) {
	authKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_AUTHKEY"))

	if err != nil {
		return nil, err
	}

	accessToken, accessExpires, err := P.accessToken(authKey, user, claims)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpires, err := P.refreshToken(authKey, user, uid, claims)

	if encryptToken {
		encKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_ENCKEY"))

		if err != nil {
			return nil, err
		}

		accessToken, err = middleware.Encrypt(encKey, accessToken)

		if err != nil {
			return nil, err
		}

		refreshToken, err = middleware.Encrypt(encKey, refreshToken)
	}

	tokens := new(Tokens)
	tokens.AccessToken = accessToken
	tokens.AccessTokenExpires = accessExpires
	tokens.RefreshToken = refreshToken
	tokens.RefreshTokenExpires = refreshExpires

	return tokens, err
}

func (P *LoginProcessor) refreshToken(key []byte, user interface{}, uid string, extraClaims map[string]interface{}) (string, int64, error) {
	expirationTime := time.Now().In(time.UTC).Add(24 * time.Hour)

	claims := make(jwt.MapClaims)

	claims["ExpiresAt"] = expirationTime.Unix()
	claims["UID"] = uid
	claims["UserType"] = P.interfaceType(user)
	claims["Extra"] = extraClaims

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString(key)

	return token, expirationTime.Unix(), err
}

func (P *LoginProcessor) accessToken(key []byte, user interface{}, extraClaims map[string]interface{}) (string, int64, error) {
	expirationTime := time.Now().In(time.UTC).Add(1 * time.Hour)

	userJson, err := json.Marshal(user)
	if err != nil {
		return "", 0, err
	}

	claims := make(jwt.MapClaims)

	claims["ExpiresAt"] = expirationTime.Unix()
	claims["User"] = userJson
	claims["UserType"] = P.interfaceType(user)
	claims["Extra"] = extraClaims

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString(key)

	return token, expirationTime.Unix(), err
}

func (P *LoginProcessor) interfaceType(x interface{}) string {
	t := reflect.TypeOf(x)

	if t.Kind() == reflect.Ptr {
		return t.Elem().String()
	}

	return t.String()
}
