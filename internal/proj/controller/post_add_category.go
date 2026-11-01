package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type AddCategory interface {
	http.Handler
}

type addCategoryController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewAddCategory(srv *server.Server, cfg data.ConfigData) AddCategory {
	return &addCategoryController{server: srv, cfg: cfg}
}

func (c *addCategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *addCategoryController) process(r *http.Request) error {
	var category data.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		return fmt.Errorf("parsing add category request: %w", err)
	}

	return c.server.CreateCategory(r.Context(), category)
}
