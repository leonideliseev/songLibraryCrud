package handler

import (
	"os"

	_ "github.com/leonideliseev/songLibraryCrud/docs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/middleware"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var validate *validator.Validate

type Handler struct {
	log *logging.Logger
	service *service.Service
}

func NewHandler(service *service.Service, log *logging.Logger) *Handler {
	log.Info("init handler...")
	return &Handler{
		log: log,
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	h.log.Info("init routes...")
	validate = validator.New()

	router := router(h.log)

	router.Use(middleware.Log(h.log))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", ping)
	api := router.Group("/api/v1")
	{
		newSongsRoutes(api.Group("/songs"), h.service.Songs, h.log)
	}

	return router
}

func router(log *logging.Logger) *gin.Engine {
	var r *gin.Engine

	if env := os.Getenv("APP_ENV"); env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
		log.Info("router for prod")
	} else {
		r = gin.Default()
		log.Info("default gin router")
	}

	return r
}
