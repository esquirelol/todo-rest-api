package get

import (
	"context"
	"errors"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/esquirelol/todo-rest-api/internal/todo"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type TaskGet interface {
	Get(ctx context.Context, author string) (todo.Todo, error)
}

func New(taskGet TaskGet, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var outTask todo.Todo
		author := chi.URLParam(r, "author")

		outTask, err := taskGet.Get(r.Context(), author)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				logger.Info("http/get:", zap.Error(storage.ErrTaskNotFound))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("task not found"))
				return
			}
			logger.Error("http/get:", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		render.JSON(w, r, outTask)

	}
}
