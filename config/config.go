package config

import (
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig(log *logging.Logger) {
	// viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error init configs: %s", err.Error())
	}
}
