package create_user

import (
	"context"
	"net/http"

	"github.com/esquirelol/todo-rest-api/internal/auth"
	"github.com/esquirelol/todo-rest-api/internal/dto"
	"github.com/esquirelol/todo-rest-api/internal/http/api/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type CreateUser interface {
	CreateUser(ctx context.Context, userName, password string) (int, error)
}

func New(create CreateUser, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userOut dto.User
		if err := render.DecodeJSON(r.Body, &userOut); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("failed to decode json")
			render.JSON(w, r, response.Error("failed to decode json"))
			return
		}

		idToken, err := create.CreateUser(r.Context(), userOut.UserName, userOut.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		token, err := auth.GenerateToken(idToken)
		if err != nil {
			logger.Error("failed to generate token", zap.Int("user_id", idToken))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to login"))
			return
		}
		logger.Info("Create user and jwt success")
		render.JSON(w, r, response.OK("jwt: "+token))
	}
}
