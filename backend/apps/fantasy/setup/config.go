package setup

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"path"
	"strings"
)

func Config() {
	conf := flag.String("config", "config.toml", "specifies config file location")
	root := flag.String("root", ".", "specifies root path for templates and public files")
	state := flag.String("state", "game.dat", "specifies file for persistent game state")

	flag.Parse()

	if conf != nil {
		configFile := path.Clean(*conf)

		switch path.Ext(configFile) {
		case "hcl":
			viper.SetConfigType("hcl")
		case "json":
			viper.SetConfigType("json")
		case "toml":
			viper.SetConfigType("toml")
		case "yaml":
			viper.SetConfigType("yaml")
		}

		viper.AddConfigPath(path.Dir(configFile))

		baseName := path.Base(configFile)

		if baseName == "." {
			baseName = "config.toml"
		}

		extName := path.Ext(baseName)
		viper.SetConfigName(baseName[:len(baseName)-len(extName)])
	} else {
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
		viper.AddConfigPath("")
	}

	// Any setting can be overridden from the environment, e.g. FANTASY_CORS_BE_NAUGHTY=true
	viper.SetEnvPrefix("FANTASY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	Defaults()

	if root != nil {
		rootpath := path.Clean(*root)
		viper.SetDefault("ROOT_PATH", rootpath)
	}

	if state != nil {
		stateFile := path.Clean(*state)
		viper.SetDefault("GAME_STATE_FILE", stateFile)
	} else {
		viper.SetDefault("GAME_STATE_FILE", "game.dat")
	}
}

func Defaults() {
	viper.SetDefault("HTTP_BASE", "/")
	viper.SetDefault("HTTP_HOST", "localhost")
	viper.SetDefault("HTTP_PORT", 80)
	viper.SetDefault("HTTP_SCHEMA", "http")

	viper.SetDefault("SITE_MODE", "dev")
	viper.SetDefault("DEBUG_USER", 0)

	viper.SetDefault("CORS_ALLOW_ORIGINS", []string{"http://localhost:5173"})
	viper.SetDefault("CORS_BE_NAUGHTY", false)
	viper.SetDefault("TRUST_PROXY", false)
	viper.SetDefault("TRUSTED_PROXIES", []string{})

	viper.SetDefault("MAX_TEAMS", 3)
	viper.SetDefault("MAX_NAME_LENGTH", 32)
	viper.SetDefault("MAX_SUBMISSION_LENGTH", 280)
	viper.SetDefault("SUBMISSION_COOLDOWN_SECONDS", 10)

	viper.SetDefault("SNAPSHOT_INTERVAL_SECONDS", 30)
}
