package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"ping": "pong",
	})
}
