package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (s *Server) MarkSessionExpired(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE session SET expired_ind = 'Y' WHERE token = ?`, token)
	if err != nil {
		return fmt.Errorf("marking session expired: %w", err)
	}

	return nil
}

func (s *Server) ExtendSession(ctx context.Context, token string, expiryDt time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE session SET expiry_dt = ? WHERE token = ?`, expiryDt, token)
	if err != nil {
		return fmt.Errorf("extending session: %w", err)
	}

	return nil
}

// GetSessionByToken returns nil, nil when no session has the given token.
func (s *Server) GetSessionByToken(ctx context.Context, token string) (*data.Session, error) {
	var session data.Session
	err := s.db.QueryRowContext(ctx,
		`SELECT session_id, token, user_id, login, name, expiry_dt, expired_ind = 'Y' FROM session WHERE token = ?`,
		token).Scan(&session.SessionId, &session.Token, &session.UserId, &session.Login, &session.Name, &session.ExpiryDt, &session.Expired)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying session by token: %w", err)
	}

	return &session, nil
}
