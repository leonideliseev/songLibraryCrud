package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	handerror "github.com/leonideliseev/songLibraryCrud/internal/handler/error"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/middleware"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	songConvert "github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

const OK = http.StatusOK

type songRouter struct {
	log     *logging.Logger
	service service.Songs
}

func newSongsRoutes(g *gin.RouterGroup, service service.Songs, log *logging.Logger) {
	log.Info("init song router...")
	r := &songRouter{
		log:     log,
		service: service,
	}

	g.GET("/", r.getSongs, middleware.CheckLimit(), middleware.CheckOffset()) // получение библиотеки с пагинацией
	g.POST("/", r.createSong)                                                 // добавление новой песни

	id := g.Group("/id", middleware.CheckId())
	{
		id.GET("", r.getSong, middleware.CheckLimit(), middleware.CheckOffset()) // получение текста песни
		id.PATCH("", r.updateSong)                                               // изменение данных песни
		id.DELETE("", r.deleteSong)                                              // удаление песни
	}
}

func (h *songRouter) getSongs(c *gin.Context) {
	limit := limitCtx(c)
	offset := offsetCtx(c)
	pagModel := &models.Song{ // модель, которая будет считывать поля фильтрации
		GroupName:   c.Query("group_name"),
		Name:        c.Query("name"),
		ReleaseDate: c.Query("release_date"),
		Text:        c.Query("text"),
		Link:        c.Query("link"),
	}

	songs, err := h.service.GetAll(c, limit, offset, pagModel)
	if err != nil {
		h.log.WithError(err).Info("get songs error")
		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseGetSongs{
		Songs: songs,
	})
}

func (h *songRouter) createSong(c *gin.Context) {
	input := &dto.RequestCreateSong{}
	if err := c.ShouldBindJSON(input); err != nil {
		h.log.WithError(err).Info("failed to read request data")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		h.log.WithError(err).Info("failed to validate request data")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songDetail, err := getSongDetailsFromAPI(input.Group, input.Name, h.log)
	if err != nil {
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	song, err := h.service.Create(c, songConvert.FromInputToModel(input, songDetail))
	if err != nil {
		if errors.Is(err, service.ErrSongAlreadyExists) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseCreateSong{
		Song: song,
	})
}

func (h *songRouter) getSong(c *gin.Context) {
	id := uuidCtx(c)
	limit := limitCtx(c)
	offset := offsetCtx(c)

	songData, err := h.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("song with id <%s> not found", id))
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	verses := strings.Split(songData.Text, "\n\n")
	if offset > len(verses) {
		handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("offset can`t be more than verses in song"))
		return
	}

	if limit > len(verses) {
		limit = len(verses) + 1
	}

	selectVerses := verses[offset:limit]
	songData.Text = strings.Join(selectVerses, "\n\n")

	c.JSON(OK, dto.ResponseGetSong{
		Song: songData,
	})
}

func (h *songRouter) updateSong(c *gin.Context) {
	id := uuidCtx(c)

	var input *dto.RequestUpdateSong
	if err := c.ShouldBindJSON(input); err != nil {
		h.log.WithError(err).Info("failed to read request data")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		h.log.WithError(err).Info("failed to validate request data")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songData, err := h.service.UpdateById(c, id, songConvert.FromInputUpdateToModel(input))
	if err != nil {
		if errors.Is(err, service.ErrUpdatedSongNotChanged) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, "song not changed")
			return
		}

		if errors.Is(err, service.ErrSongAlreadyExists) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("song with group <%s> and name <%s> already exists", songData.GroupName, songData.Name))
			return
		}

		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("song with id <%s> not found", id))
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseUpdateSong{
		Song: songData,
	})
}

func (h *songRouter) deleteSong(c *gin.Context) {
	id := uuidCtx(c)

	if err := h.service.DeleteById(c, id); err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("song with id <%s> not found", id))
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(OK)
}

func uuidCtx(c *gin.Context) uuid.UUID {
	uuidCtx, _ := c.Get(middleware.UuidCtx)
	return uuidCtx.(uuid.UUID)
}

func limitCtx(c *gin.Context) int {
	uuidCtx, _ := c.Get(middleware.LimitCtx)
	return uuidCtx.(int)
}

func offsetCtx(c *gin.Context) int {
	uuidCtx, _ := c.Get(middleware.OffsetCtx)
	return uuidCtx.(int)
}

func getSongDetailsFromAPI(group, song string, log *logging.Logger) (*dto.SongDetail, error) {
	apiURL := os.Getenv("EXTERNAL_API_URL")
	if apiURL == "" {
		log.Error("EXTERNAL_API_URL not set in environment")
		return nil, errors.New("EXTERNAL_API_URL not set in environment")
	}

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	resp, err := http.Get(apiURL + "/info" + "?" + params.Encode())
	if err != nil {
		log.WithError(err).Error("Error making GET request")
		return nil, fmt.Errorf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Get status code <%d> from external API", resp.StatusCode)
		return nil, fmt.Errorf("Error: received status code %d from external API", resp.StatusCode)
	}

	var songDetail dto.SongDetail
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Error reading response body")
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling response")
		return nil, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return &songDetail, nil
}
