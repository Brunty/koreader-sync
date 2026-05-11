package main

import (
	"net/http"

	database "github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/handlers"
	"github.com/brunty/koreader-sync-server/sync_progress"
	userpackage "github.com/brunty/koreader-sync-server/user"
)

type ServeMux struct {
	*http.ServeMux
}

func (mux ServeMux) RegisterRoutes() ServeMux {
	// See https://github.com/Open-Audiobook/koreader-sync-protocol#api-reference-summary for the reference of the
	// endpoints and spec for how this server should work

	db := database.DBCon
	mux.Handle("GET /{$}", http.HandlerFunc(handlers.Home))

	userRepo := userpackage.NewUserRepository(db)
	authMiddleware := userpackage.NewAuthMiddleware(userRepo)
	userHandler := userpackage.NewUserHandler(userRepo)
	syncHandler := sync_progress.NewSyncProgressHandler(sync_progress.NewSyncProgressRepository(db))

	mux.Handle("POST /users/create", http.HandlerFunc(userHandler.CreateUser))

	// The following routes need an auth'd user to access them
	mux.Handle("GET /users/auth", authMiddleware.Handle(http.HandlerFunc(userHandler.AuthUser)))
	mux.Handle("PUT /syncs/progress", authMiddleware.Handle(http.HandlerFunc(syncHandler.StoreSyncProgress)))
	mux.Handle("GET /syncs/progress/{document}", authMiddleware.Handle(http.HandlerFunc(syncHandler.ReadSyncProgress)))

	mux.Handle("/{path...}", http.HandlerFunc(handlers.NotFound))

	return mux
}
