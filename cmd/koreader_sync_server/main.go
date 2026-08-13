package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/brunty/koreader-sync-server/internal/db"
	"github.com/brunty/koreader-sync-server/internal/logger"
	"github.com/brunty/koreader-sync-server/internal/server"
)

func init() {
	logger.Init()

	// We use Init rather than init because we want to specify a DB file and be able to return errorsF
	err := db.Init("./data/data.db.sqlite3")
	if err != nil {
		slog.Error("database init error", slog.Any("error", err))
		return
	}

	db.SetupTables()
	slog.Debug("DB tables setup")
}

func main() {
	defer db.DBCon.Close()

	mux := server.NewServeMux()

	mux.RegisterRoutes()

	port, found := os.LookupEnv("PORT")

	if !found {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	log.Fatal(srv.ListenAndServe())
}
