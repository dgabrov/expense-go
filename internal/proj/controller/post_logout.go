package controller

import (
	"log/slog"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type Logout interface {
	http.Handler
}

type logoutController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewLogout(srv *server.Server, cfg data.ConfigData) Logout {
	return &logoutController{server: srv, cfg: cfg}
}

func (c *logoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *logoutController) process(r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	if err = c.server.MarkSessionExpired(r.Context(), cookie.Value); err != nil {
		slog.Error("logout: could not expire session", "error", err)
	}

	return nil
}
