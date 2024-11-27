package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env: %s", err.Error())
	}
}
