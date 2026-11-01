package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type UpdateCategory interface {
	http.Handler
}

type updateCategoryController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewUpdateCategory(srv *server.Server, cfg data.ConfigData) UpdateCategory {
	return &updateCategoryController{server: srv, cfg: cfg}
}

func (c *updateCategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *updateCategoryController) process(r *http.Request) error {
	var category data.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		return fmt.Errorf("parsing update category request: %w", err)
	}

	return c.server.UpdateCategory(r.Context(), category)
}
