package main

import (
	"context"

	"github.com/esquirelol/todo-rest-api/internal/config"
	"github.com/esquirelol/todo-rest-api/internal/logger"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	log := logger.LoggerInit(cfg.Env)
	log.Info("logger init", zap.String("env", cfg.Env))
	db := storage.ConnectionStorage(ctx, cfg.Storage, log)
	log.Info("successful connection to the storage")
	_ = db
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(cfg.Timeout))
}
