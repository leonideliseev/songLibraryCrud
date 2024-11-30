package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonideliseev/songLibraryCrud/config"
	"github.com/leonideliseev/songLibraryCrud/internal/handler"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	"github.com/leonideliseev/songLibraryCrud/internal/utils"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/spf13/viper"
)

// @title Song Library API
// @version 1.0
// @description API for managing a library of songs, including creating, retrieving, updating, and deleting songs.

// @contact.name Leonid Eliseev
// @contact.url https://t.me/Lenchiiiikkkk
// @contact.email leonid.2004eliseev@mail.ru

// @BasePath /api/v1
func main() {
	logger := logging.GetLogger()
	logger.Info("log writing started")
	config.InitConfig(logger)
	config.LoadEnv(logger)

	srv := initServer(logger)
	startServer(srv, logger)
	waitForShutdown(srv, logger)
}

func initServer(log *logging.Logger) *http.Server {
	var repo *repository.Repository
	utils.RepoChoice(repo, log)
	
	serv := service.NewService(repo, log)
	hand := handler.NewHandler(serv, log)

	router := hand.InitRoutes()
	return &http.Server{
		Addr:    viper.GetString("port"),
		Handler: router,
	}
}

func startServer(srv *http.Server, log *logging.Logger) {
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	log.Info("Song Library App started")
}

func waitForShutdown(srv *http.Server, log *logging.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Warn("Song Library App shutting down")

	if err := srv.Close(); err != nil {
		log.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	log.Info("Song Library App stopped")
}
