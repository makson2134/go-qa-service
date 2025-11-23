package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/makson2134/go-qa-service/internal/api"
	"github.com/makson2134/go-qa-service/internal/api/handlers"
	"github.com/makson2134/go-qa-service/internal/config"
	"github.com/makson2134/go-qa-service/internal/repository/postgres"
	"github.com/makson2134/go-qa-service/pkg"
	"github.com/pressly/goose/v3"
)

func main() {
	configPath, err := findConfigFile("config")
	if err != nil {
		log.Fatalf("failed to find config file: %v", err) // Critical error, app can't run without сonfig
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err) // Same issue here again
	}

	logger := pkg.NewLogger(cfg.Log.Level, cfg.Log.Format)

	dsn, err := cfg.Database.GetDSN()
	if err != nil {
		logger.Error("failed to get DSN", "error", err)
		log.Fatal(err)
	}

	db, err := postgres.New(
		dsn,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database connection", "error", err)
		}
	}()

	logger.Info("Connected to database")

	sqlDB, err := db.GetDB()
	if err != nil {
		logger.Error("Failed to get database instance", "error", err)
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set goose dialect", "error", err)
		log.Fatal(err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		log.Fatal(err)
	}
	logger.Info("Migrations applied successfully")

	h := handlers.New(db, db, logger)

	mux := api.SetupRoutes(h)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("starting server", "port", cfg.Server.Port, "env", cfg.Env)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-quit
	logger.Info("received shutdown signal", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info("shutting down server gracefully")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server stopped")
}
func findConfigFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", os.ErrNotExist
}
