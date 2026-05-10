package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/handlers"
	"github.com/brunty/koreader-sync-server/logger"
	"github.com/brunty/koreader-sync-server/middleware"
)

func main() {
	logger.Init()

	err := db.Init("./data/data.db.sqlite3")
	if err != nil {
		slog.Error("database init error", slog.Any("error", err))
		return
	}
	defer db.DBCon.Close()

	db.CreateTables()
	slog.Debug("DB Tables created")

	mux := http.NewServeMux()

	// See https://github.com/Open-Audiobook/koreader-sync-protocol#api-reference-summary for the reference of the
	// endpoints and spec for how this server should work

	mux.Handle("GET /{$}", http.HandlerFunc(handlers.Home))

	mux.Handle("POST /users/create", http.HandlerFunc(handlers.CreateUser))

	// The following routes need an auth'd user to access them
	mux.Handle("GET /users/auth", middleware.AuthMiddleware(http.HandlerFunc(handlers.AuthUser)))
	mux.Handle("PUT /syncs/progress", middleware.AuthMiddleware(http.HandlerFunc(handlers.StoreSyncProgress)))
	mux.Handle("GET /syncs/progress/{document}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetSyncProgress)))

	mux.Handle("/{path...}", http.HandlerFunc(handlers.NotFound))

	port, found := os.LookupEnv("PORT")

	if !found {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
