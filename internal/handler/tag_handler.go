package handler

import (
	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/pkg/apperr"
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
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodeTagList, constant.MsgInternalServerError, err))
		return
	}
	response.Success(c, data)
}
