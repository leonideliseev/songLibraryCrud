package config

import (
	"github.com/joho/godotenv"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

func LoadEnv(log *logging.Logger) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %s", err.Error())
	}

	log.Info("env loaded successfully")
}
