package main

import (
	"net/http"

	"github.com/brunty/koreader-sync-server/handlers"
	"github.com/brunty/koreader-sync-server/middleware"
)

type ServeMux struct {
	*http.ServeMux
}

func (mux ServeMux) RegisterRoutes() ServeMux {
	// See https://github.com/Open-Audiobook/koreader-sync-protocol#api-reference-summary for the reference of the
	// endpoints and spec for how this server should work

	mux.Handle("GET /{$}", http.HandlerFunc(handlers.Home))

	mux.Handle("POST /users/create", http.HandlerFunc(handlers.CreateUser))

	// The following routes need an auth'd user to access them
	mux.Handle("GET /users/auth", middleware.AuthMiddleware(http.HandlerFunc(handlers.AuthUser)))
	mux.Handle("PUT /syncs/progress", middleware.AuthMiddleware(http.HandlerFunc(handlers.StoreSyncProgress)))
	mux.Handle("GET /syncs/progress/{document}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetSyncProgress)))

	mux.Handle("/{path...}", http.HandlerFunc(handlers.NotFound))

	return mux
}
