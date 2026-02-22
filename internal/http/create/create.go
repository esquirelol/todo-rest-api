package create

import (
	"context"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/http/api/requests"
	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type TaskCreate interface {
	Create(ctx context.Context, todo requests.RequestCreate) error
}

func New(taskCr TaskCreate, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task requests.RequestCreate
		if err := render.DecodeJSON(r.Body, &task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("failed to create task")
			render.JSON(w, r, response.Error("failed to create task"))
			return
		}

		if err := taskCr.Create(r.Context(), task); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("decode json failed")
			render.JSON(w, r, response.Error("decode json failed"))
			return
		}

		logger.Info("created task success")
		render.JSON(w, r, response.OK("task is created"))
	}
}
