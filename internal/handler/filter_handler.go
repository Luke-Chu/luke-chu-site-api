package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/dto/response"
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
		response.Error(c, http.StatusInternalServerError, 50006, err.Error())
		return
	}
	response.Success(c, data)
}
