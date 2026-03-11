package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"luke-chu-site-api/internal/app/middleware"
	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/pkg/apperr"
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
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidQueryParams, constant.MsgInvalidQueryParams))
		return
	}
	data, err := h.photoService.ListPhotos(c.Request.Context(), &req)
	if err != nil {
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotosList, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) GetPhotoDetail(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidUUID, constant.MsgInvalidUUID))
		return
	}

	data, err := h.photoService.GetPhotoDetail(c.Request.Context(), photoUUID)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.ErrorFrom(c, apperr.New(404, constant.CodePhotoNotFound, constant.MsgPhotoNotFound))
			return
		}
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotoDetail, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) ViewPhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidUUID, constant.MsgInvalidUUID))
		return
	}

	visitorHash, _ := c.Get(middleware.VisitorHashKey)
	hash, _ := visitorHash.(string)
	if hash == "" {
		response.ErrorFrom(c, apperr.New(400, constant.CodeVisitorHashMissing, constant.MsgVisitorHashMissing))
		return
	}

	data, err := h.behaviorService.ViewPhoto(c.Request.Context(), photoUUID, hash)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.ErrorFrom(c, apperr.New(404, constant.CodePhotoNotFound, constant.MsgPhotoNotFound))
			return
		}
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotoDetail, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) LikePhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidUUID, constant.MsgInvalidUUID))
		return
	}

	visitorHash, _ := c.Get(middleware.VisitorHashKey)
	hash, _ := visitorHash.(string)
	if hash == "" {
		response.ErrorFrom(c, apperr.New(400, constant.CodeVisitorHashMissing, constant.MsgVisitorHashMissing))
		return
	}

	data, err := h.behaviorService.LikePhoto(c.Request.Context(), photoUUID, hash)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.ErrorFrom(c, apperr.New(404, constant.CodePhotoNotFound, constant.MsgPhotoNotFound))
			return
		}
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotoLike, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) UnlikePhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidUUID, constant.MsgInvalidUUID))
		return
	}

	visitorHash, _ := c.Get(middleware.VisitorHashKey)
	hash, _ := visitorHash.(string)
	if hash == "" {
		response.ErrorFrom(c, apperr.New(400, constant.CodeVisitorHashMissing, constant.MsgVisitorHashMissing))
		return
	}

	data, err := h.behaviorService.UnlikePhoto(c.Request.Context(), photoUUID, hash)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.ErrorFrom(c, apperr.New(404, constant.CodePhotoNotFound, constant.MsgPhotoNotFound))
			return
		}
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotoLike, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}

func (h *PhotoHandler) DownloadPhoto(c *gin.Context) {
	photoUUID := c.Param("uuid")
	if _, err := uuid.Parse(photoUUID); err != nil {
		response.ErrorFrom(c, apperr.New(400, constant.CodeInvalidUUID, constant.MsgInvalidUUID))
		return
	}

	visitorHash, _ := c.Get(middleware.VisitorHashKey)
	hash, _ := visitorHash.(string)
	if hash == "" {
		response.ErrorFrom(c, apperr.New(400, constant.CodeVisitorHashMissing, constant.MsgVisitorHashMissing))
		return
	}

	data, err := h.behaviorService.DownloadPhoto(c.Request.Context(), photoUUID, hash)
	if err != nil {
		if errors.Is(err, service.ErrPhotoNotFound) {
			response.ErrorFrom(c, apperr.New(404, constant.CodePhotoNotFound, constant.MsgPhotoNotFound))
			return
		}
		response.ErrorFrom(c, apperr.Wrap(500, constant.CodePhotoDownload, constant.MsgInternalServerError, err))
		return
	}

	response.Success(c, data)
}
