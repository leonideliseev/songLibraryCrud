package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonideliseev/songLibraryCrud/pkg/handler"
	"github.com/leonideliseev/songLibraryCrud/pkg/repository"
	"github.com/leonideliseev/songLibraryCrud/pkg/repository/postgres"
	"github.com/leonideliseev/songLibraryCrud/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error init configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env: %s", err.Error())
	}

	var repos *repository.Repository
	repos, err := postgres.NewPostgresRepository(postgres.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})
	
	if err != nil {
		logrus.Fatalf("failed init db: %s", err.Error())
	}

	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := gin.Default()
	handlers.InitRoutes(srv)

	go func() {
		if err := srv.Run(viper.GetString("port")); err != nil {
			logrus.Fatalf("error running server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	srv.Run(":8080")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
