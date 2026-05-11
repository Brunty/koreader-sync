package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunty/koreader-sync-server/types"
)

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&types.ErrorResponse{Error: message})
}

func WriteStatusResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: message})
}
