package delete

import (
	"context"
	"errors"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type TaskDelete interface {
	Delete(ctx context.Context, idTask string) error
}

func New(taskDelete TaskDelete, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idTask := chi.URLParam(r, "id")

		if err := taskDelete.Delete(r.Context(), idTask); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				logger.Info("http/delete: task not exists")
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("task not exists"))
				return
			}
			logger.Error("http/delete:", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		logger.Info("task is delete", zap.String("id", idTask))
		render.JSON(w, r, response.OK("task is delete"))
	}
}
