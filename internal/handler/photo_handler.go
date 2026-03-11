package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"luke-chu-site-api/internal/app/middleware"
	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/service"
)

type PhotoHandler struct {
	photoService    service.PhotoService
	behaviorService service.BehaviorService
	validate        *validator.Validate
}

func NewPhotoHandler(photoService service.PhotoService, behaviorService service.BehaviorService, validate *validator.Validate) *PhotoHandler {
	return &PhotoHandler{
		photoService:    photoService,
		behaviorService: behaviorService,
		validate:        validate,
	}
}

func (h *PhotoHandler) ListPhotos(c *gin.Context) {
	var req request.PhotoListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid query params")
		return
	}
	data, err := h.photoService.ListPhotos(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) GetPhotoDetail(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid uuid")
		return
	}

	data, err := h.photoService.GetPhotoDetail(c.Request.Context(), photoUUID)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "photo not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50002, "internal server error")
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) ViewPhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid uuid")
		return
	}

	data, err := h.behaviorService.ViewPhoto(c.Request.Context(), photoUUID)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "photo not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50002, "internal server error")
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) LikePhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid uuid")
		return
	}

	visitorHash, _ := c.Get(middleware.VisitorHashKey)
	hash, _ := visitorHash.(string)
	if hash == "" {
		response.Error(c, http.StatusBadRequest, 40003, "visitor hash missing")
		return
	}

	data, err := h.behaviorService.LikePhoto(c.Request.Context(), photoUUID, hash)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "photo not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50003, "internal server error")
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) DownloadPhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if err := h.behaviorService.DownloadPhoto(c.Request.Context(), photoUUID); err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "photo not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50004, err.Error())
		return
	}

	response.Success(c, gin.H{"uuid": photoUUID, "action": "download"})
}
