package logger

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/brunty/koreader-sync-server/request_id"
	"github.com/stretchr/testify/assert"
)

func TestLogRequestDetails(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := slog.Default()
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	handler := LogRequestDetails(nextHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "request received")
	assert.Contains(t, logOutput, "requestID")
	assert.Contains(t, logOutput, "127.0.0.1:12345")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test/path")
}

func TestLogRequestDetails_WithRequestID(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := slog.Default()
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	ctx := req.Context()
	ctx = context.WithValue(ctx, request_id.ContextKeyRequestID, "request-id-here")
	req = req.WithContext(ctx)

	handler := LogRequestDetails(nextHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "request received")
	assert.Contains(t, logOutput, "requestID=request-id-here")
	assert.Contains(t, logOutput, "127.0.0.1:12345")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test/path")
}
