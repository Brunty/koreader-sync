package request_id

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// An unexported zero-size struct type means no other package can construct this key, so nothing outside
// this package can collide with (or overwrite) the value we store under it
type requestIDContextKey struct{}

func ContextWithRequestID(r *http.Request, requestID uuid.UUID) context.Context {
	return context.WithValue(r.Context(), requestIDContextKey{}, requestID.String())
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}
