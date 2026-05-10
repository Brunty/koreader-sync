package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/types"
)

func GetSyncProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	document := r.PathValue("document")
	if document == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "not found"})
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := dao.SelectProgress(userId, document)

	if err != nil {
		slog.Debug("get sync progress error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	response, err := progress.MarshalToResponse()
	if err != nil {
		slog.Debug("get sync progress marshaling error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "something went wrong"})
		return
	}

	json.NewEncoder(w).Encode(response)

	return
}
