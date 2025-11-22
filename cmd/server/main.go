package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/makson2134/go-qa-service/internal/api"
	"github.com/makson2134/go-qa-service/internal/api/handlers"
	"github.com/makson2134/go-qa-service/internal/config"
	"github.com/makson2134/go-qa-service/internal/repository/postgres"
	"github.com/makson2134/go-qa-service/pkg"
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
		log.Fatal(err) //
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

	logger.Info("connected to database")

	h := handlers.New(db, db, logger)

	mux := api.SetupRoutes(h)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Info("starting server", "port", cfg.Server.Port, "env", cfg.Env)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
		log.Fatal(err)
	}
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
