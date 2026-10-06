package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/karbasia/doc-manager/internal/auth"
	"github.com/karbasia/doc-manager/internal/config"
	"github.com/karbasia/doc-manager/internal/database"
	"github.com/karbasia/doc-manager/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadAppConfig()
	logger := newLogger(cfg)

	// Start up and migrate DB
	db, err := database.Open(ctx, cfg.DB)
	if err != nil {
		logger.Error("open database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(ctx, db.Write, logger); err != nil {
		logger.Error("migrate database", "err", err)
		os.Exit(1)
	}
	logger.Info("database ready", "path", cfg.DB.Path)

	tokens := auth.NewTokenManager(cfg)

	srv := server.New(cfg.HTTP, server.NewRouter(logger, tokens))

	logger.Info("server starting", "addr", cfg.HTTP.Addr(), "env", cfg.Env)
	if err := srv.Run(ctx); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}

func newLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	_ = level.UnmarshalText([]byte(cfg.LogLevel)) // already validated in config

	opts := &slog.HandlerOptions{Level: level}
	if cfg.IsProduction() {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
