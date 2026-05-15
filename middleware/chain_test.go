package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_MiddlewareProcessedInOrderAdded(t *testing.T) {
	var processingOutput []string

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwA")
			next.ServeHTTP(w, r)
		})
	}
	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwB")
			next.ServeHTTP(w, r)
		})
	}
	mwC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwC")
			next.ServeHTTP(w, r)
		})
	}
	chain := Chain{mwA, mwB, mwC}

	req := httptest.NewRequest(http.MethodGet, "/some-url", nil)
	w := httptest.NewRecorder()

	chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(w, req)

	assert.Equal(t, []string{"mwA", "mwB", "mwC"}, processingOutput)
}

func TestChain_MiddlewareChainCanBeExtended(t *testing.T) {
	var processingOutput []string

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwA")
			next.ServeHTTP(w, r)
		})
	}
	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwB")
			next.ServeHTTP(w, r)
		})
	}
	mwC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwC")
			next.ServeHTTP(w, r)
		})
	}
	mwD := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processingOutput = append(processingOutput, "mwD")
			next.ServeHTTP(w, r)
		})
	}
	chain := Chain{mwA, mwB}
	newChain := chain.Extend(mwC, mwD)

	req := httptest.NewRequest(http.MethodGet, "/some-url", nil)
	w := httptest.NewRecorder()

	newChain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(w, req)

	assert.Equal(t, []string{"mwA", "mwB", "mwC", "mwD"}, processingOutput)
}
