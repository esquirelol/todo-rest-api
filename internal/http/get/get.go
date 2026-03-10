package get

import (
	"context"
	"errors"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/auth"
	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/esquirelol/todo-rest-api/internal/models"
	"github.com/esquirelol/todo-rest-api/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type TaskGet interface {
	Get(ctx context.Context, author string, userId int) ([]models.ModelTodo, error)
}

func New(taskGet TaskGet, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var outTask []models.ModelTodo
		author := chi.URLParam(r, "author")
		authHeader := r.Header.Get("Authorization")
		tokenString := authHeader[7:]
		userId, err := auth.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		outTask, err = taskGet.Get(r.Context(), author, userId)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				logger.Info("http/get:", zap.Error(storage.ErrNotExists))
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
