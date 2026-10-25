package proj

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/controller"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

const serverReadHeaderTimeout = 10 * time.Second

func newHandler(db *sql.DB, cfg data.ConfigData) http.Handler {
	mux := http.NewServeMux()
	context := cfg.Context

	srv := server.NewServer(db, &cfg)

	mux.Handle("GET "+context+"/{$}", controller.NewRoot(db, cfg))
	mux.Handle("POST "+context+"/login", controller.NewLogin(srv, cfg))
	mux.Handle("POST "+context+"/version", controller.NewVersion(srv, cfg))
	mux.Handle("POST "+context+"/getAccounts", controller.NewGetAccounts(srv, cfg))
	mux.Handle("POST "+context+"/addAccount", controller.NewAddAccount(srv, cfg))
	mux.Handle("POST "+context+"/updateAccount", controller.NewUpdateAccount(srv, cfg))
	mux.Handle("POST "+context+"/deleteAccounts", controller.NewDeleteAccounts(srv, cfg))
	mux.Handle("POST "+context+"/getGroups", controller.NewGetGroups(srv, cfg))
	mux.Handle("POST "+context+"/addGroup", controller.NewAddGroup(srv, cfg))
	mux.Handle("POST "+context+"/updateGroup", controller.NewUpdateGroup(srv, cfg))
	mux.Handle("POST "+context+"/deleteGroups", controller.NewDeleteGroups(srv, cfg))
	mux.Handle("POST "+context+"/getCategories", controller.NewGetCategories(srv, cfg))
	mux.Handle("POST "+context+"/addCategory", controller.NewAddCategory(srv, cfg))
	mux.Handle("POST "+context+"/updateCategory", controller.NewUpdateCategory(srv, cfg))
	mux.Handle("POST "+context+"/deleteCategories", controller.NewDeleteCategories(srv, cfg))
	mux.Handle("POST "+context+"/deleteTransactions", controller.NewDeleteTransactions(srv, cfg))
	mux.Handle("POST "+context+"/getCompleteGroups", controller.NewGetCompleteGroups(srv, cfg))
	mux.Handle("POST "+context+"/getTransactionById", controller.NewGetTransactionById(srv, cfg))
	mux.Handle("POST "+context+"/searchTransactions", controller.NewSearchTransactions(srv, cfg))
	mux.Handle("POST "+context+"/addTransaction", controller.NewAddTransaction(srv, cfg))
	mux.Handle("POST "+context+"/updateTransaction", controller.NewUpdateTransaction(srv, cfg))
	mux.Handle("POST "+context+"/logout", controller.NewLogout(srv, cfg))

	return mux
}

func serve(cfg data.ConfigData, db *sql.DB) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(db, cfg),
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	slog.Info("http server listening", "addr", addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("serving http on %s: %w", addr, err)
	}

	return nil
}
