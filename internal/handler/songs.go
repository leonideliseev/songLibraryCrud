package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	"github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/sirupsen/logrus"
)

const OK = http.StatusOK

type songRouter struct {
	service service.Songs
}

func newSongsRoutes(g *gin.RouterGroup, service service.Songs) {
	r := &songRouter{
		service: service,
	}

	g.GET("/", r.getSongs)  // получение библиотеки с пагинацией
	g.POST("/", r.createSong)  // добавление новой песни
	g.GET("/:id", r.getSong)  // получение текста песни
	g.PATCH("/:id", r.updateSong)  // изменение данных песни
	g.DELETE("/:id", r.deleteSong)  // удаление песни
}

func (h *songRouter) getSongs(c *gin.Context) {
	limit := getDefaultQuery(c, "limit", "10")
    offset := getDefaultQuery(c, "offset", "0")

	songs, err := h.service.GetAll(limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseGetSongs{
		Songs: songs,
	})
}

func (h *songRouter) createSong(c *gin.Context) {
	var input dto.RequestCreateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	song, err := h.service.CreateSong(songConvert.FromHandToServ(input))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseCreateSong{
		Song: &song,
	})
}

func (h *songRouter) getSong(c *gin.Context) {
	group, song := getGroupAndSong(c)

	songData, err := h.service.GetSong(group, song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseGetSong{
		Song: &songData,
	})
}

// TODO: сделать получение данных с помощью dto.UpdateSong
func (h *songRouter) updateSong(c *gin.Context) {
	group, song := getGroupAndSong(c)

	updatedData := models.Song{
		GroupName: group,
		Name: song,
		ReleaseDate: "ppp",
		Text: "ppp",
		Link: "ppp",
	}

	songData, err := h.service.UpdateSong(group, song, updatedData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseUpdateSong{
		Song: &songData,
	})
}

func (h *songRouter) deleteSong(c *gin.Context) {
	group, song := getGroupAndSong(c)

	if err := h.service.DeleteSong(group, song); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(OK)
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func getDefaultQuery(c *gin.Context, name, def string) int {
	param := c.DefaultQuery(name, def)

	intParam, err := strconv.Atoi(param)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 0
	}

	if intParam < 0 {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s can't be negative", name))
		return 0
	}

	return intParam
}

func getGroupAndSong(c *gin.Context) (string, string) {
	group := c.Query("group")
    song := c.Query("song")

	if group == "" || song == "" {
		newErrorResponse(c, http.StatusBadRequest, "group and song parameters are required")
		return "", ""
	}

	return group, song
}
