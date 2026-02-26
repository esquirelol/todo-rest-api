package main

import (
	"context"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/config"
	"github.com/esquirelol/todo-rest-api/internal/http/create"
	del "github.com/esquirelol/todo-rest-api/internal/http/delete"
	"github.com/esquirelol/todo-rest-api/internal/http/get"
	"github.com/esquirelol/todo-rest-api/internal/http/update"
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
	db, err := storage.ConnectionStorage(ctx, cfg.Storage, log)
	if err != nil {
		log.Error("failed to connect storage", zap.Error(err))
		return
	}
	log.Info("successful connection to the storage")
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(cfg.Timeout))
	router.Post("/create", create.New(&db, log))
	router.Get("/{author}", get.New(&db, log))
	router.Patch("/{id}", update.New(&db, log))
	router.Delete("/{id}", del.New(&db, log))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	srv.ListenAndServe()
}
