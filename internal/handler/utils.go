package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/middleware"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

const OK = http.StatusOK

var uuidPath = fmt.Sprintf("/:%s", middleware.UuidCtx)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Error reading response body")
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	var songDetail dto.SongDetail
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling response")
		return nil, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return &songDetail, nil
}
