package config

import (
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/spf13/viper"
)

func InitConfig(log *logging.Logger) {
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error init configs: %s", err.Error())
	}

	log.Info("config readed")
}
