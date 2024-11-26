package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
)

func InitRoutes(router *gin.Engine, service *service.Service) {
	router.GET("/ping", ping)
	api := router.Group("/api/v1")
	{
		newSongsRoutes(api.Group("/songs"), service.Songs)
	}
}
