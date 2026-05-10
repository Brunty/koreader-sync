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
		writeErrorResponse(w, http.StatusBadRequest, "username and password are required")
		return
	}

	if err := req.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := req.MarshalToUser()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	err = dao.StoreUser(user)
	if err != nil {
		slog.Error("store user error", slog.Any("error", err))

		// I hate having to do string matching here, try to find a better way to match (error code if possible?)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			writeErrorResponse(w, http.StatusBadRequest, "username is already taken")
			return
		}

		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	writeStatusResponse(w, http.StatusCreated, "user created")
	return
}
