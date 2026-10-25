package proj

import (
	"database/sql"
	"log/slog"
)

func Start() (*sql.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	slog.Info("config loaded", "version", cfg.Version)

	db, err := openDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	slog.Info("database connected", "host", cfg.Database.Host, "port", cfg.Database.Port, "name", cfg.Database.Name)

	return db, nil
}
