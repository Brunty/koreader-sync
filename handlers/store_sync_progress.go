package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/types"
)

func StoreSyncProgress(w http.ResponseWriter, r *http.Request) {
	req := &types.SyncProgressRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Debug("store sync progress bad request body", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "bad body content"})
		return
	}

	if err := req.Validate(); err != nil {
		slog.Debug("create user validation failed", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: err.Error()})
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := req.MarshalToProgress(userId)
	if err != nil {
		slog.Debug("store sync marshaling failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	// try finding an entry by user ID and document, if so, update it, if not, create it
	err = dao.StoreProgress(progress)
	if err != nil {
		slog.Error("store progress error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: "sync stored"})

	return
}
