package handler

import (
	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/dto/response"
)

type HealthHandler struct {
	serviceName string
}

func NewHealthHandler(serviceName string) *HealthHandler {
	return &HealthHandler{serviceName: serviceName}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"service": h.serviceName,
	})
}
