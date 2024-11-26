package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
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

	g.GET("/", r.getSongs)
	g.POST("/", r.addSong)
	g.GET("/:id", r.getSong)
	g.PUT("/:id", r.updateSong)
	g.DELETE("/id", r.deleteSong)
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

func (h *songRouter) getSongs(c *gin.Context) {
	limit := getDefaultQuery(c, "limit", "10")
    offset := getDefaultQuery(c, "offset", "0")

	songs, err := h.service.GetAll(limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, songs)
}

func (h *songRouter) addSong(c *gin.Context) {
	var input models.Song
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	song, err := h.service.CreateSong(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, song)
}

func (h *songRouter) getSong(c *gin.Context) {
	group, song := getGroupAndSong(c)

	songData, err := h.service.GetSong(group, song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, songData)
}

func (h *songRouter) updateSong(c *gin.Context) {
	group, song := getGroupAndSong(c)

	updatedData := &models.Song{
		Group: group,
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

	c.JSON(OK, songData)
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