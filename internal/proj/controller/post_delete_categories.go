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
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		return fmt.Errorf("parsing delete categories request: %w", err)
	}

	// TODO
	return nil
}
