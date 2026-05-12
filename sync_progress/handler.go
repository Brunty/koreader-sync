package sync_progress

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/handlers"
)

type SyncProgressHandler struct {
	repo SyncProgressRepository
}

func NewSyncProgressHandler(repository SyncProgressRepository) *SyncProgressHandler {
	return &SyncProgressHandler{repo: repository}
}

func (h *SyncProgressHandler) ReadSyncProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	document := r.PathValue("document")
	if document == "" {
		handlers.WriteErrorResponse(w, http.StatusNotFound, "not found")
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := h.repo.SelectByUserIDAndDocument(userId, document)

	if err != nil {
		slog.Debug("get sync progress error", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	if progress == nil {
		slog.Debug("progress not found")
		handlers.WriteErrorResponse(w, http.StatusNotFound, "not found")
		return
	}

	response, err := progress.MarshalToReadResponse()
	if err != nil {
		slog.Debug("get sync progress marshaling error", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (h *SyncProgressHandler) StoreSyncProgress(w http.ResponseWriter, r *http.Request) {
	req := &StoreSyncProgressRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Debug("store sync progress bad request body", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusBadRequest, "bad body content")
		return
	}

	if err := req.Validate(); err != nil {
		slog.Debug("create user validation failed", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := r.Context().Value("user").(int64)
	progress, err := req.MarshalToSyncProgress(userId)
	if err != nil {
		slog.Debug("store sync marshaling failed", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	_, err = h.repo.Store(progress)
	if err != nil {
		slog.Error("store progress error", slog.Any("error", err))
		handlers.WriteErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	handlers.WriteStatusResponse(w, http.StatusOK, "sync stored")
}
