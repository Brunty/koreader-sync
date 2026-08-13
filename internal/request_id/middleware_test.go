package request_id

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItAddsRequestId(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEqual(t, "", r.Context().Value(ContextKeyRequestID))
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	handler := AddRequestIDToMiddleware(nextHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}
