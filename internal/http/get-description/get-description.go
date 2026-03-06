package get_description

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

type TaskGetDescription interface {
	GetDescription(ctx context.Context, idTask string) (models.ModelDescription, error)
}

func New(taskGD TaskGetDescription, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var outModel models.ModelDescription
		idTask := chi.URLParam(r, "id")

		outModel, err := taskGD.GetDescription(r.Context(), idTask)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.OK("Task dont exists"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return

		}
		render.JSON(w, r, outModel)
	}
}
