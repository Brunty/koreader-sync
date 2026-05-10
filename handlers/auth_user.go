package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunty/koreader-sync-server/types"
)

func AuthUser(w http.ResponseWriter, _ *http.Request) {
	// We don't needto do any auth checking here because this handler is protected by middleware.AuthMiddleware and
	// so if we've got here, we're auth'd, so let's just return a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: "authorized"})
}
