package controller

import (
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetGroups interface {
	http.Handler
}

type getGroupsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetGroups(srv *server.Server, cfg data.ConfigData) GetGroups {
	return &getGroupsController{server: srv, cfg: cfg}
}

func (c *getGroupsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getGroupsController) process(r *http.Request) ([]data.Group, error) {
	// TODO
	return nil, nil
}
