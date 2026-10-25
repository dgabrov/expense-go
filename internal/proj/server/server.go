package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dgabrov/expense-go/internal/data"
)

type Server struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServer(db *sql.DB, config *data.ConfigData) *Server {
	return &Server{db, config}
}

func (s *Server) CreateSession(ctx context.Context, session data.Session) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO session (session_id, token, user_id, login, name, expiry_dt) VALUES (?, ?, ?, ?, ?, ?)`,
		session.SessionId, session.Token, session.UserId, session.Login, session.Name, session.ExpiryDt)
	if err != nil {
		return fmt.Errorf("inserting session: %w", err)
	}

	return nil
}
