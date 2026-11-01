package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type DeleteCategories interface {
	http.Handler
}

type deleteCategoriesController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewDeleteCategories(srv *server.Server, cfg data.ConfigData) DeleteCategories {
	return &deleteCategoriesController{server: srv, cfg: cfg}
}

func (c *deleteCategoriesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *deleteCategoriesController) process(r *http.Request) error {
	var holder data.IdsHolder
	if err := json.NewDecoder(r.Body).Decode(&holder); err != nil {
		return fmt.Errorf("parsing delete categories request: %w", err)
	}

	return c.server.DeleteCategories(r.Context(), holder.Ids)
}
