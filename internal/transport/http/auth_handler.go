package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/DENFNC/web-test/internal/transport/dto/request"
	"github.com/DENFNC/web-test/internal/transport/dto/response"
	"github.com/DENFNC/web-test/internal/utils/mapping"
)

type AuthService interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	LoginUser(ctx context.Context, user *domain.User) (string, error)
	RevokeToken(ctx context.Context, token string) error
}

type AuthHandler struct {
	*slog.Logger
	AuthService
}

func NewAuthHandler(log *slog.Logger, mux *http.ServeMux, srv AuthService) {
	handler := &AuthHandler{
		Logger:      log,
		AuthService: srv,
	}

	{
		mux.HandleFunc("POST /api/register", handler.register)
		mux.HandleFunc("POST /api/auth", handler.auth)
		mux.HandleFunc("DELETE /api/auth/{token}", handler.exit)
	}
}

func (api *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterUserRequest
	if !decodeAndValidate(w, r, &req) {
		return
	}

	var user domain.User
	if err := mapping.MapStruct(req, &user); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to map request")
		return
	}

	login, err := api.AuthService.CreateUser(r.Context(), &user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "User creation failed")
		return
	}

	response.JSON(w, http.StatusOK, response.RegisterUserResponse{Login: login})
}

func (api *AuthHandler) auth(w http.ResponseWriter, r *http.Request) {
	var req request.AuthUserRequest
	if !decodeAndValidate(w, r, &req) {
		return
	}

	var user domain.User
	if err := mapping.MapStruct(req, &user); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to map request")
		return
	}

	token, err := api.AuthService.LoginUser(r.Context(), &user)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Login failed")
		return
	}

	response.JSON(w, http.StatusOK, response.AuthUserResponse{Token: token})
}

func (api *AuthHandler) exit(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		response.Error(w, http.StatusBadRequest, "missing token")
		return
	}

	err := api.AuthService.RevokeToken(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "cannot revoke token")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"response": map[string]bool{
			token: true,
		},
	})
}

func decodeAndValidate[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return false
	}

	if v, ok := any(dst).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			response.Error(w, http.StatusBadRequest, "Validation failed")
			return false
		}
	}

	return true
}
