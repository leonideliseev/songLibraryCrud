package app

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
	"github.com/spf13/viper"
)

type Closer interface {
	Close()
}

type App struct {
	logger *logging.Logger

	srv *http.Server
	conn postgresql.Conn

	repo *repository.Repository
	serv *service.Service
	hand *handler.Handler

	quit chan os.Signal
}

func NewApp() *App {
	app := &App{}

	app.logger = logging.GetLogger()
	app.logger.Info("log writing started")

	config.InitConfig(app.logger)
	config.LoadEnv(app.logger)

	app.initDBConn()
	app.initAppCore()
	app.initServer()
	app.initShutdown()

	return app
}

func (a *App) Run() {
	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			a.logger.Fatalf("error running server: %s", err.Error())
		}
	}()

	a.logger.Info("Song Library App started")

	<-a.quit
	a.logger.Warn("Song Library App shutting down")

	if err := a.srv.Close(); err != nil {
		a.logger.Errorf("error occurred on server shutting down: %s", err.Error())
	}
	a.conn.Close()

	a.logger.Info("Song Library App stopped")
}

func (a *App) initAppCore() {
	a.repo = repository.New(a.conn, a.logger)
	a.serv = service.New(a.repo, a.logger)
	a.hand = handler.New(a.serv, a.logger)
}

func (a *App) initServer() {
	router := a.hand.InitRoutes()

	a.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", viper.GetString("http.host"), viper.GetString("http.port")),
		Handler: router,
	}
}

func (a *App) initShutdown() {
	a.quit = make(chan os.Signal, 1)
	signal.Notify(a.quit, syscall.SIGTERM, syscall.SIGINT)
}
