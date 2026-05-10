package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunty/koreader-sync-server/types"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&types.ErrorResponse{Error: "not found"})
}
