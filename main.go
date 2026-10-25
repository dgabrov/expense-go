package main

import (
	"github.com/dgabrov/expense-go/internal/proj"
	"log/slog"
)

func main() {
	db, err := proj.Start()

	if err != nil {
		slog.Error("failed to start project", "error", err)
		return
	}
	defer db.Close()
}
