package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	handerror "github.com/leonideliseev/songLibraryCrud/internal/handler/error"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/middleware"
	"github.com/leonideliseev/songLibraryCrud/internal/service"
	"github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
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

	id := g.Group("/id", middleware.CheckId())
	{
		id.GET("", r.getSong)  // получение текста песни
		id.PATCH("", r.updateSong)  // изменение данных песни
		id.DELETE("", r.deleteSong)  // удаление песни
	}
}

// TODO: сделать фильтры по полям
func (h *songRouter) getSongs(c *gin.Context) {
	limit := getDefaultQuery(c, "limit", "10")
    offset := getDefaultQuery(c, "offset", "0")

	songs, err := h.service.GetAll(c, limit, offset)
	if err != nil {
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
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songDetail, err := getSongDetailsFromAPI(input.Group, input.Name)
	if err != nil {
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	song, err := h.service.CreateSong(c, songConvert.FromInputToModel(input, songDetail))
	if err != nil {
		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseCreateSong{
		Song: song,
	})
}

func (h *songRouter) getSong(c *gin.Context) {
	id := uuidCtx(c)
	limit := getDefaultQuery(c, "limit", "100")
    offset := getDefaultQuery(c, "offset", "0")

	songData, err := h.service.GetSong(c, id)
	if err != nil {
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
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songData, err := h.service.UpdateSong(c, id, songConvert.FromInputUpdateToModel(input))
	if err != nil {
		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(OK, dto.ResponseUpdateSong{
		Song: songData,
	})
}

func (h *songRouter) deleteSong(c *gin.Context) {
	id := uuidCtx(c)

	if err := h.service.DeleteSong(c, id); err != nil {
		handerror.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(OK)
}

func getDefaultQuery(c *gin.Context, name, def string) int {
	param := c.DefaultQuery(name, def)

	intParam, err := strconv.Atoi(param)
	if err != nil {
		handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return 0
	}

	if intParam < 0 {
		handerror.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s can't be negative", name))
		return 0
	}

	return intParam
}

func uuidCtx(c *gin.Context) uuid.UUID {
	uuidCtx, _ := c.Get(middleware.UuidCtx)
	return uuidCtx.(uuid.UUID)
}

func getSongDetailsFromAPI(group, song string) (*dto.SongDetail, error) {
	apiURL := os.Getenv("EXTERNAL_API_URL") // TODO: изменить
	if apiURL == "" {
		return nil, fmt.Errorf("EXTERNAL_API_URL not set in environment")
	}

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: received status code %d from external API", resp.StatusCode)
	}

	var songDetail dto.SongDetail
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return &songDetail, nil
}
