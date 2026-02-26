package update

import (
	"context"
	"errors"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/dto"
	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type TaskUpdate interface {
	Update(ctx context.Context, todo dto.TodoUpdate, idTask string) error
}

func New(taskUpdate TaskUpdate, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo dto.TodoUpdate
		idTask := chi.URLParam(r, "id")
		if err := render.DecodeJSON(r.Body, &todo); err != nil {
			logger.Error("failed to decode task", zap.String("id", idTask), zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := taskUpdate.Update(r.Context(), todo, idTask); err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				logger.Info("http/update: task dont found")
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("task not found"))
				return
			}
			logger.Error("http/update:", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		logger.Info("task is update", zap.String("id", idTask))
		render.JSON(w, r, response.OK("task is update"))
	}
}
