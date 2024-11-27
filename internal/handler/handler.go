package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
)

var validate *validator.Validate

func InitRoutes(router *gin.Engine, service *service.Service) {
	validate = validator.New()
	
	router.GET("/ping", ping)
	api := router.Group("/api/v1")
	{
		newSongsRoutes(api.Group("/songs"), service.Songs)
	}
}
