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

	mux := &ServeMux{http.NewServeMux()}

	mux.RegisterRoutes()

	port, found := os.LookupEnv("PORT")

	if !found {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
