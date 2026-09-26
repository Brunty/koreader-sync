package request_id

import (
	"net/http"

	"github.com/google/uuid"
)

func AddRequestIDToMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()

		r = r.WithContext(ContextWithRequestID(r, requestID))

		next.ServeHTTP(w, r)
	})
}
