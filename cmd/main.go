package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonideliseev/songLibraryCrud/config"
	"github.com/leonideliseev/songLibraryCrud/internal/handler"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
	"github.com/leonideliseev/songLibraryCrud/schema"
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

	srv, conn := initServer(logger)
	startServer(srv, logger)
	waitForShutdown(srv, conn, logger)
}

func initServer(log *logging.Logger) (*http.Server, Closer) {
	config := postgresql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	configTest := config
	configTest.DBName = "postgres"
	connTest, err := postgresql.ConnWithPgxPool(configTest)
	if err != nil {
		log.Fatal("")
	}
	defer connTest.Close()

	postgresql.CreateDatabaseIfNotExists(connTest, viper.GetString("db.dbname"), log)

	conn, err := postgresql.ConnWithPgxPool(config)
	if err != nil {
		log.Fatal("")
	}

	postgresql.Migrate(log, &schema.DB, &config)

	repo := repository.New(conn, log)
	serv := service.New(repo, log)
	hand := handler.New(serv, log)

	router := hand.InitRoutes()
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString("port")),
		Handler: router,
	}, conn
}

func startServer(srv *http.Server, log *logging.Logger) {
	defer log.Info("Song Library App started")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()
}

func waitForShutdown(srv *http.Server, conn Closer, log *logging.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Warn("Song Library App shutting down")

	if err := srv.Close(); err != nil {
		log.Errorf("error occurred on server shutting down: %s", err.Error())
	}
	conn.Close()

	log.Info("Song Library App stopped")
}

type Closer interface {
	Close()
}
