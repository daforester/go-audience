package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"io"
	"time"
)

func ValidLoginToken(data string, decryptToken bool) (*jwt.Token, error) {
	authKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_AUTHKEY"))

	if err != nil {
		return nil, err
	}

	if decryptToken {
		encKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_ENCKEY"))

		if err != nil {
			return nil, err
		}

		data, err = Decrypt(encKey, data)
		if err != nil {
			return nil, err
		}
	}

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return authKey, nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	expiry, ok := token.Claims.(jwt.MapClaims)["ExpiresAt"]

	if !ok {
		return nil, errors.New("no expiry set")
	} else if int64(expiry.(float64)) < time.Now().In(time.UTC).Unix() {
		return nil, errors.New("token expired")
	}

	return token, nil
}

func Encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func Decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipher-text block size is too short")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}
