package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"nexus/internal/adapter/http/v1/router"
	"nexus/internal/infrastructure/config"
	"nexus/internal/infrastructure/database"
	"nexus/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	jwtpkg "nexus/pkg/jwt"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Configuration loaded",
		slog.String("environment", cfg.App.Environment),
		slog.String("version", cfg.App.Version))

	// Connect to db
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", slog.Any("error", err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", slog.Any("error", err))
		}
	}()

	// Init JWT
	jwtManager := jwtpkg.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// Init modules
	healthRouter := router.InitHealthModule()

	// HTTP server

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Routes
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		v1Router := router.NewV1Router(healthRouter)
		v1Router.Setup(v1)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		host := cfg.Server.Host
		if host == "0.0.0.0" || host == "" {
			host = "localhost"
		}

		logger.Info("Server started",
			slog.Int("port", cfg.Server.Port),
			slog.String("host", cfg.Server.Host),
			slog.String("health_check", fmt.Sprintf("http://%s:%d/api/v1/health", host, cfg.Server.Port)),
			slog.String("environment", cfg.App.Environment),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", slog.Any("error", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Exiting server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to exit", slog.Any("error", err))
	}

	logger.Info("Server exited gracefully")
}
