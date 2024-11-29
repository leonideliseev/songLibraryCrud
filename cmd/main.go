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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := logging.GetLogger()
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
			logrus.Fatalf("error running server: %s", err.Error())
		}
	}()
	logrus.Print("MainApp started")
}

func waitForShutdown(srv *http.Server, log *logging.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Close(); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
