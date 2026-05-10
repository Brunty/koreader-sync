package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/types"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	req := &types.CreateUserRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Debug("create user bad request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "username and password are required"})
		return
	}

	if err := req.Validate(); err != nil {
		slog.Debug("create user validation failed", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := req.MarshalToUser()
	if err != nil {
		slog.Debug("create user marshaling failed", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	err = dao.StoreUser(user)
	if err != nil {
		slog.Error("store user error", slog.Any("error", err))

		// I hate having to do string matching here, try to find a better way to match (error code if possible?)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "username is already taken"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: "user created"})
	return
}
