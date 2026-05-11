package handlers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	writeStatusResponse(w, http.StatusOK, "This is a server compatible with the KOReader Sync Protocol")
	return
}
