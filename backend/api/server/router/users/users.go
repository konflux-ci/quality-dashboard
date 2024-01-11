package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	userV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/user/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
	"go.uber.org/zap"
)

func (s *usersRouter) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var user userV1Alpha1.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading user value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	err := s.Storage.CreateUser(user)
	fmt.Println("error creating user:", err)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully created user",
		StatusCode: http.StatusCreated,
	})
}

func (s *usersRouter) getUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	userEmail := r.URL.Query()["user_email"]

	if len(userEmail) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "user_email value not present in query",
			StatusCode: 400,
		})
	}

	user, err := s.Storage.GetUser(userEmail[0])
	if err != nil {
		s.Logger.Error("Failed to fetch user. Make sure the user exists", zap.String("user", userEmail[0]), zap.Error(err))

		// if returns error, dex session will expire
		// return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
		// 	Message:    err.Error(),
		// 	StatusCode: http.StatusBadRequest,
		// })
	}

	return httputils.WriteJSON(w, http.StatusOK, user)
}

func (s *usersRouter) listAllUsers(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	users, err := s.Storage.ListAllUsers()
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	fmt.Println("aqui vai users", users)

	return httputils.WriteJSON(w, http.StatusOK, users)
}
