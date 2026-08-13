package handlers

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

func WriteStatusResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&StatusResponse{Status: message})
}
