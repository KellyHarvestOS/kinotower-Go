package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KellyHarvestOS/kinotower-Go/internal/config"
	"github.com/KellyHarvestOS/kinotower-Go/internal/database"
	"github.com/KellyHarvestOS/kinotower-Go/internal/handlers"
	"github.com/KellyHarvestOS/kinotower-Go/internal/middleware"
	"github.com/KellyHarvestOS/kinotower-Go/internal/repositories"
	"github.com/KellyHarvestOS/kinotower-Go/internal/routes"
	"github.com/KellyHarvestOS/kinotower-Go/internal/server"
	"github.com/KellyHarvestOS/kinotower-Go/internal/services"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx := context.Background()

	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if cfg.AutoMigrate {
		if err := database.RunMigrations(db, "migrations"); err != nil {
			log.Error("migrations failed", "error", err)
			os.Exit(1)
		}
		log.Info("migrations applied")
	}
	if cfg.MigrateOnly {
		return
	}

	if cfg.AutoSeed {
		if err := database.RunSeeds(ctx, db, "seeds"); err != nil {
			log.Error("seeds failed", "error", err)
			os.Exit(1)
		}
		log.Info("seeds applied")
	}
	if cfg.SeedOnly {
		return
	}

	repo := repositories.New(db)
	service := services.New(repo, cfg.JWTSecret, cfg.JWTTTL)
	handler := handlers.New(service)
	authMiddleware := middleware.Auth(cfg.JWTSecret)
	router := routes.New(handler, authMiddleware)

	appHandler := middleware.Chain(
		router,
		middleware.Recover(log),
		middleware.RequestID,
		middleware.Logger(log),
		middleware.CORS,
		middleware.JSONContent,
		middleware.RateLimit(cfg.RateLimitRPM),
		middleware.Timeout(cfg.RequestTimeout),
	)

	srv := server.New(cfg, appHandler, log)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	case sig := <-stop:
		log.Info("shutdown signal received", "signal", sig.String())
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}
