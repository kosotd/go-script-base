package logger

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func Instance() *logrus.Logger {
	return logrus.StandardLogger()
}
