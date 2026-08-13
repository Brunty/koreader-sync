package middleware

import (
	"net/http"
	"slices"
)

type Chain []func(http.Handler) http.Handler

func (c Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.Then(h)
}

// Then Middleware is processed in the order added to the chain
func (c Chain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}

func (c Chain) Extend(handlers ...func(http.Handler) http.Handler) Chain {
	return append(c, handlers...)
}
