package main

import (
	"log/slog"
	"nexus/internal/infrastructure/config"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Configuration loaded",
		slog.String("environment", cfg.App.Environment),
		slog.String("version", cfg.App.Version))
}
