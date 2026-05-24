package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/logger"
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

	mux := &ServeMux{http.NewServeMux()}

	mux.RegisterRoutes()

	port, found := os.LookupEnv("PORT")

	if !found {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
