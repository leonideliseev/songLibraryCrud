package config

import (
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/spf13/viper"
)

func InitConfig(log *logging.Logger) {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error init configs: %s", err.Error())
	}

	log.Info("config readed successfully")
}
