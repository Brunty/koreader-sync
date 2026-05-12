package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/brunty/koreader-sync-server/handlers"
)

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) AuthUser(w http.ResponseWriter, _ *http.Request) {
	// We don't needto do any auth checking here because this handler is protected by middleware.AuthMiddleware and
	// so if we've got here, we're auth'd, so let's just return a success message
	w.Header().Set("Content-Type", "application/json")
	handlers.WriteStatusResponse(w, http.StatusOK, "authorized")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := &CreateUserRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handlers.WriteErrorResponse(w, http.StatusBadRequest, "username and password are required")
		return
	}

	if err := req.Validate(); err != nil {
		handlers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := req.MarshalToUser()
	if err != nil {
		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	_, err = h.repo.Store(user)
	if err != nil {
		slog.Error("store user error", slog.Any("error", err))

		// I hate having to do string matching here, try to find a better way to match (error code if possible?)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			handlers.WriteErrorResponse(w, http.StatusBadRequest, "username is already taken")
			return
		}

		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	handlers.WriteStatusResponse(w, http.StatusCreated, "user created")
}
