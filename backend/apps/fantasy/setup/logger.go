package setup

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Logger() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	switch viper.GetInt("LOG_LEVEL") {
	case 1: // ERROR
		logrus.SetLevel(logrus.ErrorLevel)
	case 2: // WARN
		logrus.SetLevel(logrus.WarnLevel)
	case 3: // INFO
		logrus.SetLevel(logrus.InfoLevel)
	case 4: // DEBUG
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.FatalLevel)
	}
}
