package processors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
	ValidLoginToken(data string, decryptToken bool) (*jwt.Token, error)
	LoginTokens(user interface{}, uid string, encryptToken bool, claims map[string]interface{}) (*Tokens, error)
}

type LoginProcessor struct {
}

func (P LoginProcessor) New() *LoginProcessor {
	return new(LoginProcessor)
}

func (P *LoginProcessor) ValidLoginToken(data string, decryptToken bool) (*jwt.Token, error) {
	authKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_AUTHKEY"))

	if err != nil {
		return nil, err
	}

	if decryptToken {
		encKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_ENCKEY"))

		if err != nil {
			return nil, err
		}

		data, err = P.decrypt(encKey, data)
		if err != nil {
			return nil, err
		}
	}

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return authKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	expiry, ok := token.Claims.(jwt.MapClaims)["ExpiresAt"]

	if !ok {
		return nil, errors.New("no expiry set")
	}

	expiresAt, ok := expiry.(float64)
	if !ok {
		return nil, errors.New("invalid expiry")
	} else if int64(expiresAt) < time.Now().In(time.UTC).Unix() {
		return nil, errors.New("token expired")
	}

	return token, nil
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
	if err != nil {
		return nil, err
	}

	if encryptToken {
		encKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_ENCKEY"))

		if err != nil {
			return nil, err
		}

		accessToken, err = P.encrypt(encKey, accessToken)

		if err != nil {
			return nil, err
		}

		refreshToken, err = P.encrypt(encKey, refreshToken)

		if err != nil {
			return nil, err
		}
	}

	tokens := new(Tokens)
	tokens.AccessToken = accessToken
	tokens.AccessTokenExpires = accessExpires
	tokens.RefreshToken = refreshToken
	tokens.RefreshTokenExpires = refreshExpires

	return tokens, nil
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
	expirationTime := time.Now().In(time.UTC).Add(5 * time.Minute)

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

// encrypt seals the message with AES-GCM, the nonce is prefixed to the cipher-text
func (P *LoginProcessor) encrypt(key []byte, message string) (encmess string, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(message), nil)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func (P *LoginProcessor) decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	if len(cipherText) < gcm.NonceSize() {
		err = errors.New("cipher-text is too short")
		return
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return
	}

	decodedmess = string(plainText)
	return
}
