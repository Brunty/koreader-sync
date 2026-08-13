package logger

import (
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/internal/handlers"
	"github.com/brunty/koreader-sync-server/internal/request_id"
)

func LogRequestDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		requestID := r.Context().Value(request_id.ContextKeyRequestID)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)

		slog.Info("request received", slog.Any("requestID", requestID), userAttrs, requestAttrs)
		next.ServeHTTP(w, r)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", slog.Any("panic", r))
				w.Header().Set("Content-Type", "application/json")
				handlers.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
