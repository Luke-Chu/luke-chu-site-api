package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"luke-chu-site-api/internal/app/middleware"
	"luke-chu-site-api/internal/handler"
)

func NewRouter(
	logger *zap.Logger,
	behaviorGuardCfg middleware.BehaviorGuardConfig,
	healthHandler *handler.HealthHandler,
	photoHandler *handler.PhotoHandler,
	tagHandler *handler.TagHandler,
	filterHandler *handler.FilterHandler,
) *gin.Engine {
	engine := gin.New()

	engine.Use(middleware.CORS())
	engine.Use(middleware.Logger(logger))
	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.Visitor())

	v1 := engine.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)

		v1.GET("/photos", photoHandler.ListPhotos)
		v1.GET("/photos/:uuid", photoHandler.GetPhotoDetail)

		v1.GET("/tags", tagHandler.ListTags)
		v1.GET("/filters", filterHandler.GetFilters)
	}

	behaviorGroup := v1.Group("/photos/:uuid")
	behaviorGroup.Use(middleware.BehaviorGuard(logger, behaviorGuardCfg))
	{
		behaviorGroup.POST("/view", photoHandler.ViewPhoto)
		behaviorGroup.POST("/like", photoHandler.LikePhoto)
		behaviorGroup.POST("/unlike", photoHandler.UnlikePhoto)
		behaviorGroup.POST("/download", photoHandler.DownloadPhoto)
	}

	return engine
}
