package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type DeleteGroups interface {
	http.Handler
}

type deleteGroupsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewDeleteGroups(srv *server.Server, cfg data.ConfigData) DeleteGroups {
	return &deleteGroupsController{server: srv, cfg: cfg}
}

func (c *deleteGroupsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *deleteGroupsController) process(r *http.Request) error {
	var holder data.IdsHolder
	if err := json.NewDecoder(r.Body).Decode(&holder); err != nil {
		return fmt.Errorf("parsing delete groups request: %w", err)
	}

	return c.server.DeleteGroups(r.Context(), holder.Ids)
}
