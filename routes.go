package main

import (
	"net/http"

	"github.com/brunty/koreader-sync-server/auth"
	database "github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/handlers"
	"github.com/brunty/koreader-sync-server/logger"
	"github.com/brunty/koreader-sync-server/middleware"
	"github.com/brunty/koreader-sync-server/request_id"
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

	userRepo := userpackage.NewUserRepository(db)
	authMiddleware := auth.NewAuthMiddleware(userRepo)

	userHandler := userpackage.NewUserHandler(userRepo)
	syncHandler := sync_progress.NewSyncProgressHandler(sync_progress.NewSyncProgressRepository(db))

	// Middleware is processed in the order they are added to the chain
	baseMiddlewareChain := middleware.Chain{request_id.AddRequestIDToMiddleware, logger.LogRequestDetails}
	authMiddlewareChain := baseMiddlewareChain.Extend(authMiddleware.Handle)

	// The following routes don't need an auth'd user to access them
	mux.Handle("GET /{$}", baseMiddlewareChain.ThenFunc(handlers.Home))
	mux.Handle("/{path...}", baseMiddlewareChain.ThenFunc(handlers.NotFound))
	mux.Handle("POST /users/create", baseMiddlewareChain.ThenFunc(userHandler.CreateUser))

	// The following routes need an auth'd user to access them
	mux.Handle("GET /users/auth", authMiddlewareChain.ThenFunc(userHandler.AuthUser))
	mux.Handle("PUT /syncs/progress", authMiddlewareChain.ThenFunc(syncHandler.StoreSyncProgress))
	mux.Handle("GET /syncs/progress/{document}", authMiddlewareChain.ThenFunc(syncHandler.ReadSyncProgress))

	return mux
}
