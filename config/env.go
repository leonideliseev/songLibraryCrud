package config

import (
	"github.com/joho/godotenv"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/sirupsen/logrus"
)

func LoadEnv(log *logging.Logger) {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env: %s", err.Error())
	}
}
