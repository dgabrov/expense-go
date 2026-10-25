package proj

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/go-sql-driver/mysql"
)

const dbPingTimeout = 5 * time.Second

func openDB(cfg data.DatabaseConfig) (*sql.DB, error) {
	dsn := mysql.NewConfig()
	dsn.User = cfg.User
	dsn.Passwd = cfg.Password
	dsn.Net = "tcp"
	dsn.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	dsn.DBName = cfg.Name
	dsn.ParseTime = true

	connector, err := mysql.NewConnector(dsn)
	if err != nil {
		return nil, fmt.Errorf("building database connector: %w", err)
	}

	db := sql.OpenDB(connector)

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connecting to database %s: %w", dsn.Addr, err)
	}

	return db, nil
}
