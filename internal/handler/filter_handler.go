package handler

import (
	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/pkg/apperr"
	"luke-chu-site-api/internal/service"
)

type FilterHandler struct {
	filterService service.FilterService
}

func NewFilterHandler(filterService service.FilterService) *FilterHandler {
	return &FilterHandler{filterService: filterService}
}

func (h *FilterHandler) GetFilters(c *gin.Context) {
	data, err := h.filterService.GetFilters(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodeFilterList, constant.MsgInternalServerError, err))
		return
	}
	response.Success(c, data)
}
