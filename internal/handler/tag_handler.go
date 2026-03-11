package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/service"
)

type TagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

func (h *TagHandler) ListTags(c *gin.Context) {
	data, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50005, err.Error())
		return
	}
	response.Success(c, data)
}
