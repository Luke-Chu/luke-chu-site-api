package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/pkg/apperr"
)

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, CommonResponse{
		Code:    constant.CodeSuccess,
		Message: constant.MsgSuccess,
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

func ErrorFrom(c *gin.Context, err error) {
	var appErr *apperr.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
		return
	}

	Error(c, http.StatusInternalServerError, constant.CodeInternalServer, constant.MsgInternalServerError)
}
