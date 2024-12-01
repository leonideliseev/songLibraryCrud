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
	"time"

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

var uuidPath = fmt.Sprintf("/:%s", middleware.UuidCtx)

type songRouter struct {
	log     *logging.Logger
	service service.Songs
}

func newSongsRoutes(g *gin.RouterGroup, service service.Songs, log *logging.Logger) {
	log.Info("init song routes...")
	r := &songRouter{
		log:     log,
		service: service,
	}

	g.GET("", middleware.CheckLimit(), middleware.CheckOffset(), r.getSongs) // получение библиотеки с пагинацией
	g.POST("", r.createSong)                                                 // добавление новой песни

	id := g.Group(uuidPath, middleware.CheckId())
	{
		id.GET("", middleware.CheckLimit(), middleware.CheckOffset(), r.getSong) // получение текста песни
		id.PATCH("", r.updateSong)                                               // изменение данных песни
		id.DELETE("", r.deleteSong)                                              // удаление песни
	}
}

// @Summary Get Songs
// @Description Retrieves a list of songs with optional filters for group name, name, release date, text, and link.
// @Tags songs
// @Produce json
// @Param group_name query string false "Filter by group name"
// @Param name query string false "Filter by song name"
// @Param release_date query string false "Filter by release date (format: YYYY-MM-DD)"
// @Param text query string false "Filter by text in the song"
// @Param link query string false "Filter by link"
// @Param limit query int false "Maximum number of items to retrieve (pagination)"
// @Param offset query int false "Number of items to skip (pagination)"
// @Success 200 {object} dto.ResponseGetSongs "Successful response with a list of songs"
// @Failure 500 {object} handerror.ErrorResponse "Internal server error"
// @Router /songs [get]
func (h *songRouter) getSongs(c *gin.Context) {
	limit := limitCtx(c)
	offset := offsetCtx(c)
	timeInput, err := timeFromQuery(c)
	if err != nil {
		h.log.WithError(err).Info("failed to validate time")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pagModel := &models.Song{ // модель, которая будет считывать поля фильтрации
		GroupName:   c.Query("group_name"),
		Name:        c.Query("name"),
		ReleaseDate: timeInput,
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

// @Summary Create Song
// @Description Creates a new song with the provided details.
// @Tags songs
// @Accept json
// @Produce json
// @Param input body dto.RequestCreateSong true "Details of the song to be created"
// @Success 201 {object} dto.ResponseCreateSong "Song created successfully"
// @Failure 400 {object} handerror.ErrorResponse "Bad request: validation error or song already exists"
// @Failure 409 {object} handerror.ErrorResponse "Conflict: song with the specified group and name already exists"
// @Failure 500 {object} handerror.ErrorResponse "Internal server error"
// @Router /songs [post]
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

	convSong, err := songConvert.FromInputToModel(input, songDetail)
	if err != nil {
		h.log.WithError(err).Info("failed to validate time")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	song, err := h.service.Create(c, convSong)
	if err != nil {
		if errors.Is(err, service.ErrSongAlreadyExists) {
			handerror.NewErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseCreateSong{
		Song: song,
	})
}

// @Summary Get Song
// @Description Retrieves a song by its ID. Supports pagination for song text verses.
// @Tags songs
// @Produce json
// @Param uuid path string true "Song ID (validated as UUID)"
// @Param limit query int false "Maximum number of verses to retrieve (pagination)"
// @Param offset query int false "Number of verses to skip (pagination)"
// @Success 200 {object} dto.ResponseGetSong "Successful response with the song details"
// @Failure 400 {object} handerror.ErrorResponse "Bad request: invalid ID, invalid limit/offset, or song not found"
// @Failure 404 {object} handerror.ErrorResponse "Not found: song with the specified ID does not exist"
// @Failure 500 {object} handerror.ErrorResponse "Internal server error"
// @Router /songs/{uuid} [get]
func (h *songRouter) getSong(c *gin.Context) {
	id := uuidCtx(c)
	limit := limitCtx(c)
	offset := offsetCtx(c)

	songData, err := h.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusNotFound, fmt.Sprintf("song with id <%s> not found", id))
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

// @Summary Update Song
// @Description Updates a song by its ID with new details.
// @Tags songs
// @Accept json
// @Produce json
// @Param uuid path string true "Song ID (validated as UUID)"
// @Param input body dto.RequestUpdateSong true "Details for updating the song"
// @Success 200 {object} dto.ResponseUpdateSong "Successful response with updated song details"
// @Failure 400 {object} handerror.ErrorResponse "Bad request: invalid id, input or song not changed"
// @Failure 404 {object} handerror.ErrorResponse "Not found: song with the specified ID does not exist"
// @Failure 409 {object} handerror.ErrorResponse "Conflict: song with the specified group and name already exists"
// @Failure 500 {object} handerror.ErrorResponse "Internal server error"
// @Router /songs/{uuid} [patch]
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

	convSong, err := songConvert.FromInputUpdateToModel(input)
	if err != nil {
		h.log.WithError(err).Info("failed to validate time")
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songData, err := h.service.UpdateById(c, id, convSong)
	if err != nil {
		if errors.Is(err, service.ErrUpdatedSongNotChanged) {
			handerror.NewErrorResponse(c, http.StatusBadRequest, "song not changed")
			return
		}

		if errors.Is(err, service.ErrSongAlreadyExists) {
			handerror.NewErrorResponse(c, http.StatusConflict, fmt.Sprintf("song with group <%s> and name <%s> already exists", songData.GroupName, songData.Name))
			return
		}

		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusNotFound, fmt.Sprintf("song with id <%s> not found", id))
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseUpdateSong{
		Song: songData,
	})
}

// @Summary Delete Song
// @Description Deletes a song by its ID.
// @Tags songs
// @Param uuid path string true "Song ID (validated as UUID)"
// @Success 204 "Song successfully deleted"
// @Failure 400 {object} handerror.ErrorResponse "Bad request: invalid ID format"
// @Failure 404 {object} handerror.ErrorResponse "Not found: song with the specified ID does not exist"
// @Failure 500 {object} handerror.ErrorResponse "Internal server error"
// @Router /songs/{uuid} [delete]
func (h *songRouter) deleteSong(c *gin.Context) {
	id := uuidCtx(c)

	if err := h.service.DeleteById(c, id); err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			handerror.NewErrorResponse(c, http.StatusNotFound, fmt.Sprintf("song with id <%s> not found", id))
			return
		}

		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func timeFromQuery(c *gin.Context) (time.Time, error) {
	t := c.Query("release_date")

	if t == "" {
		return time.Time{}, nil
	}

	return time.Parse("2006-01-02", t)
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
