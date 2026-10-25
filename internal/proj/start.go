package proj

import (
	"log/slog"
)

func Start() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	slog.Info("config loaded", "version", cfg.Version)

	db, err := openDB(cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("database connected", "host", cfg.Database.Host, "port", cfg.Database.Port, "name", cfg.Database.Name)

	return serve(cfg, db)
}
