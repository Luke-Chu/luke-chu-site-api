package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"luke-chu-site-api/internal/config"
)

func NewServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
}
