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
	Delete(ctx context.Context, title string) error
}

func New(taskDelete TaskDelete, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := chi.URLParam(r, "title")

		if err := taskDelete.Delete(r.Context(), title); err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				logger.Info("http/delete: task not found")
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("task not found"))
				return
			}
			logger.Error("http/delete:", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		logger.Info("task is delete", zap.String("title", title))
		render.JSON(w, r, response.OK("task is delete"))
	}
}
