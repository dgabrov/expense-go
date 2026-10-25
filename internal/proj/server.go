package proj

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/controller"
)

const serverReadHeaderTimeout = 10 * time.Second

func newHandler(db *sql.DB, context string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+context+"/{$}", controller.GetRoot)

	return mux
}

func serve(cfg data.ConfigData, db *sql.DB) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(db, cfg.Context),
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	slog.Info("http server listening", "addr", addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("serving http on %s: %w", addr, err)
	}

	return nil
}
