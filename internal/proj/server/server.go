package server

import (
	"database/sql"

	"github.com/dgabrov/expense-go/internal/data"
)

type Server struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServer(db *sql.DB, config *data.ConfigData) *Server {
	return &Server{db, config}
}
