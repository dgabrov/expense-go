package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type AddGroup interface {
	http.Handler
}

type addGroupController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewAddGroup(srv *server.Server, cfg data.ConfigData) AddGroup {
	return &addGroupController{server: srv, cfg: cfg}
}

func (c *addGroupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *addGroupController) process(r *http.Request) error {
	var group data.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		return fmt.Errorf("parsing add group request: %w", err)
	}

	return c.server.CreateGroup(r.Context(), group)
}
