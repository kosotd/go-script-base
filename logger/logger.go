package logger

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("log_level", "INFO", "log level")

	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
