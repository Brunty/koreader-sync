package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel_Debug(t *testing.T) {
	t.Setenv("LOG_LEVEL", "DEBUG")
	assert.Equal(t, slog.LevelDebug, getLogLevel())
}

func TestGetLogLevel_Info(t *testing.T) {
	t.Setenv("LOG_LEVEL", "INFO")
	assert.Equal(t, slog.LevelInfo, getLogLevel())
}

func TestGetLogLevel_Warn(t *testing.T) {
	t.Setenv("LOG_LEVEL", "WARN")
	assert.Equal(t, slog.LevelWarn, getLogLevel())
}

func TestGetLogLevel_Error(t *testing.T) {
	t.Setenv("LOG_LEVEL", "ERROR")
	assert.Equal(t, slog.LevelError, getLogLevel())
}

func TestGetLogLevel_Default(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	assert.Equal(t, slog.LevelWarn, getLogLevel())
}

func TestGetLogLevel_UnknownValue(t *testing.T) {
	t.Setenv("LOG_LEVEL", "TRACE")
	assert.Equal(t, slog.LevelWarn, getLogLevel())
}

func TestGetLogLevel_Unset(t *testing.T) {
	assert.Equal(t, slog.LevelWarn, getLogLevel())
}

func TestInit(t *testing.T) {
	Logger = nil

	Init()

	assert.NotNil(t, Logger)
	assert.Equal(t, slog.Default(), Logger)
}
