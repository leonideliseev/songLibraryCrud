package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

func Log(log *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(map[string]interface{}{
			"method": c.Request.Method,
			"uri": c.Request.URL.Path},
			).Debug("request")

		c.Next()
	}
}
