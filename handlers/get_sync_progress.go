package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/dao"
)

func GetSyncProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	document := r.PathValue("document")
	if document == "" {
		writeErrorResponse(w, http.StatusNotFound, "not found")
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := dao.SelectProgress(userId, document)

	if err != nil {
		slog.Debug("get sync progress error", slog.Any("error", err))
		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	response, err := progress.MarshalToResponse()
	if err != nil {
		slog.Debug("get sync progress marshaling error", slog.Any("error", err))
		writeErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	json.NewEncoder(w).Encode(response)

	return
}
