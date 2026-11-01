package controller

import (
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type Version interface {
	http.Handler
}

type versionController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewVersion(srv *server.Server, cfg data.ConfigData) Version {
	return &versionController{server: srv, cfg: cfg}
}

func (c *versionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *versionController) process(r *http.Request) (string, error) {
	return c.cfg.Version, nil
}
