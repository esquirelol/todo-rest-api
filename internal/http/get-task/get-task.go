package get_task

import (
	"context"
	"errors"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/esquirelol/todo-rest-api/internal/models"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type GetTask interface {
	GetId(ctx context.Context, idTask string) (models.ModelTodo, error)
}

func New(getTask GetTask, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idTask := chi.URLParam(r, "id")

		task, err := getTask.GetId(r.Context(), idTask)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				w.WriteHeader(http.StatusNotFound)
				logger.Info("task dont exists", zap.String("id-task", idTask))
				response.OK("task dont exists")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("internal server error", zap.String("id-task", idTask), zap.Error(err))
			response.Error("internal sever error")
			return
		}
		render.JSON(w, r, task)
		logger.Info("get task is success")
	}
}
