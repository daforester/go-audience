package setup

import (
	"bitbucket.org/daforester/go-backend-repository/component/libs/strgen"
	"encoding/base32"
	"fmt"
	"github.com/spf13/viper"
)

func Keys() {
	if viper.GetString("SESSION_AUTHKEY") == "" {
		authKey, err := strgen.RandBytes(64)
		if err != nil {
			panic(fmt.Errorf("failed to generate random bytes for auth key"))
		}

		viper.SetDefault("SESSION_AUTHKEY", base32.StdEncoding.EncodeToString(authKey))
		viper.Set("SESSION_AUTHKEY", base32.StdEncoding.EncodeToString(authKey))
		fmt.Println("GENERATED AUTHENTICATION KEY")
		fmt.Println(viper.GetString("SESSION_AUTHKEY"))
	}
	if viper.GetString("SESSION_ENCKEY") == "" {
		encKey, err := strgen.RandBytes(32)
		if err != nil {
			panic(fmt.Errorf("failed to generate random bytes for encryption key"))
		}

		viper.SetDefault("SESSION_ENCKEY", base32.StdEncoding.EncodeToString(encKey))
		viper.Set("SESSION_ENCKEY", base32.StdEncoding.EncodeToString(encKey))
		fmt.Println("GENERATED ENCRYPTION KEY")
		fmt.Println(viper.GetString("SESSION_ENCKEY"))
	}
}
