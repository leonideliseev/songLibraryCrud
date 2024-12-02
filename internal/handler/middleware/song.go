package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/handerr"
)

const (
	UuidCtx = "uuid"
	LimitCtx = "limit"
	OffsetCtx = "offset"
)

func CheckId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(UuidCtx)

		uuid, err := uuid.Parse(id)
		if err != nil {
			handerr.NewErrorResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		c.Set(UuidCtx, uuid)

		c.Next()
	}
}

func CheckLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		readedParam := c.DefaultQuery(LimitCtx, "10")

		num, err := strconv.Atoi(readedParam)
		if err != nil {
			handerr.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if num < 1 {
			handerr.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s can't be less than 1: <%d>", LimitCtx, num))
			return
		}

		c.Set(LimitCtx, num)

		c.Next()
	}
}

func CheckOffset() gin.HandlerFunc {
	return func(c *gin.Context) {
		readedParam := c.DefaultQuery(OffsetCtx, "0")

		num, err := strconv.Atoi(readedParam)
		if err != nil {
			handerr.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if num < 0 {
			handerr.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s can't be less than 0: <%d>", OffsetCtx, num))
			return
		}


		c.Set(OffsetCtx, num)

		c.Next()
	}
}
