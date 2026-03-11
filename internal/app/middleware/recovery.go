package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	dtoresponse "luke-chu-site-api/internal/dto/response"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					zap.Any("panic", rec),
					zap.String("path", c.Request.URL.Path),
					zap.ByteString("stack", debug.Stack()),
				)

				dtoresponse.Error(c, http.StatusInternalServerError, 50000, fmt.Sprintf("internal server error"))
				c.Abort()
			}
		}()

		c.Next()
	}
}
