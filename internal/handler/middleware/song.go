package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handerror "github.com/leonideliseev/songLibraryCrud/internal/handler/error"
)

const (
	UuidCtx = "uuid"
)

func CheckId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		uuid, err := uuid.Parse(id)
		if err != nil {
			handerror.NewErrorResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		c.Set(UuidCtx, uuid)

		c.Next()
	}
}
