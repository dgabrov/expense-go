package controller

import (
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
	// TODO
	return nil
}
