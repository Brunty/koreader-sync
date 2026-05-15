package request_id

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const ContextKeyRequestID string = "requestID"

func AddRequestIDToMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyRequestID, requestID.String())
		r = r.WithContext(ctx)

		slog.Info("request ID attached", slog.String("requestID", requestID.String()))
		next.ServeHTTP(w, r)
	})
}
