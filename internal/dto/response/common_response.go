package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, CommonResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, CommonResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
