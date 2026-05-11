package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/not-found-url-does-not-matter-here-though", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotFound)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	expectedRsp := &ErrorResponse{Error: "not found"}
	actualRsp := &ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}
