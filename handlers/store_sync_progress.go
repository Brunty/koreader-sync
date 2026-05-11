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
		slog.Debug("store sync progress bad request body", slog.Any("error", err))
		writeErrorResponse(w, http.StatusBadRequest, "bad body content")
		return
	}

	if err := req.Validate(); err != nil {
		slog.Debug("create user validation failed", slog.Any("error", err))
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := req.MarshalToProgress(userId)
	if err != nil {
		slog.Debug("store sync marshaling failed", slog.Any("error", err))
		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// try finding an entry by user ID and document, if so, update it, if not, create it
	err = dao.StoreProgress(progress)
	if err != nil {
		slog.Error("store progress error", slog.Any("error", err))
		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	writeStatusResponse(w, http.StatusOK, "sync stored")

	return
}
