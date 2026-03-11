package main

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"luke-chu-site-api/internal/app"
	"luke-chu-site-api/internal/app/middleware"
	"luke-chu-site-api/internal/config"
	"luke-chu-site-api/internal/db"
	"luke-chu-site-api/internal/handler"
	ossutil "luke-chu-site-api/internal/pkg/oss"
	"luke-chu-site-api/internal/repository"
	"luke-chu-site-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := middleware.NewZapLogger(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	pg, err := db.NewPostgres(cfg)
	if err != nil {
		logger.Fatal("failed to initialize postgres", zap.Error(err))
	}
	defer pg.Close()

	validate := validator.New()

	photoRepo := repository.NewPhotoRepository(pg)
	tagRepo := repository.NewTagRepository(pg)
	filterRepo := repository.NewFilterRepository(pg)
	downloadSigner, err := ossutil.NewPresignDownloadURLSigner(cfg.OSS)
	if err != nil {
		logger.Fatal("failed to initialize oss signer", zap.Error(err))
	}

	photoService := service.NewPhotoService(photoRepo)
	behaviorService := service.NewBehaviorService(photoRepo, downloadSigner)
	tagService := service.NewTagService(tagRepo)
	filterService := service.NewFilterService(filterRepo)

	healthHandler := handler.NewHealthHandler(cfg.App.Name)
	photoHandler := handler.NewPhotoHandler(photoService, behaviorService, validate)
	tagHandler := handler.NewTagHandler(tagService)
	filterHandler := handler.NewFilterHandler(filterService)

	router := app.NewRouter(logger, healthHandler, photoHandler, tagHandler, filterHandler)
	server := app.NewServer(cfg, router)

	logger.Info("server starting",
		zap.String("name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("address", server.Addr),
	)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("server exited unexpectedly", zap.Error(err))
	}
}
