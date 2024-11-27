package handler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
)

var validate *validator.Validate

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	validate = validator.New()
	
	router.GET("/ping", ping)
	api := router.Group("/api/v1")
	{
		newSongsRoutes(api.Group("/songs"), h.service.Songs)
	}

	return router
}

func router() *gin.Engine {
	var r *gin.Engine

	if env := os.Getenv("APP_ENV"); env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	} else {
		r = gin.Default()
	}

	return r
}
