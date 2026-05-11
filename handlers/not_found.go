package handlers

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	writeErrorResponse(w, http.StatusNotFound, "not found")
}
