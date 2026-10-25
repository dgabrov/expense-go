package controller

import (
	"errors"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetCategories interface {
	http.Handler
}

type getCategoriesController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetCategories(srv *server.Server, cfg data.ConfigData) GetCategories {
	return &getCategoriesController{server: srv, cfg: cfg}
}

func (c *getCategoriesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getCategoriesController) process(r *http.Request) ([]data.Category, error) {
	groupId := r.URL.Query().Get("groupId")
	if groupId == "" {
		return nil, errors.New("groupId query parameter is required")
	}

	return c.server.GetCategories(r.Context(), groupId)
}
