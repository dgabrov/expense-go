package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type UpdateGroup interface {
	http.Handler
}

type updateGroupController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewUpdateGroup(srv *server.Server, cfg data.ConfigData) UpdateGroup {
	return &updateGroupController{server: srv, cfg: cfg}
}

func (c *updateGroupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *updateGroupController) process(r *http.Request) error {
	var group data.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		return fmt.Errorf("parsing update group request: %w", err)
	}

	// TODO
	return nil
}
