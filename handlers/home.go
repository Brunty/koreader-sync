package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunty/koreader-sync-server/types"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: "This is a server compatible with the KOReader Sync Protocol"})
	return
}
