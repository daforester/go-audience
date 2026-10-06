package setup

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/libs/strgen"
	"encoding/base32"
	"fmt"
	"github.com/spf13/viper"
)

func Keys() {
	// Keys are never committed to config files for Fantasy. In production they must be
	// provided (FANTASY_SESSION_AUTHKEY / FANTASY_SESSION_ENCKEY), otherwise every restart
	// would invalidate all player tokens.
	if viper.GetString("SITE_MODE") == "production" {
		if viper.GetString("SESSION_AUTHKEY") == "" || viper.GetString("SESSION_ENCKEY") == "" {
			panic(fmt.Errorf("SESSION_AUTHKEY and SESSION_ENCKEY must be set in production"))
		}
	}

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
