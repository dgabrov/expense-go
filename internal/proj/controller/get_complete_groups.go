package controller

import (
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetCompleteGroups interface {
	http.Handler
}

type getCompleteGroupsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetCompleteGroups(srv *server.Server, cfg data.ConfigData) GetCompleteGroups {
	return &getCompleteGroupsController{server: srv, cfg: cfg}
}

func (c *getCompleteGroupsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getCompleteGroupsController) process(r *http.Request) ([]data.CompleteGroup, error) {
	return c.server.GetCompleteGroups(r.Context())
}
