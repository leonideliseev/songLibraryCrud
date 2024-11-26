package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/songLibraryCrud/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h Handler) InitRoutes(s *gin.Engine) {
	api := s.Group("/api")
	{
		api.GET("/ping", h.ping)

		songs := api.Group("songs")
		{
			songs.GET("/", h.getSongs)
			songs.POST("/", h.addSong)
			songs.GET("/:id", h.getSong)
			songs.PUT("/:id", h.updateSong)
			songs.DELETE("/id", h.deleteSong)
		}
	}
}