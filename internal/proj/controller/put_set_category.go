package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type SetCategory interface {
	http.Handler
}

type setCategoryController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewSetCategory(srv *server.Server, cfg data.ConfigData) SetCategory {
	return &setCategoryController{server: srv, cfg: cfg}
}

func (c *setCategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *setCategoryController) process(r *http.Request) error {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return err
	}

	var request data.SetCategoryData
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return fmt.Errorf("parsing set category request: %w", err)
	}

	if request.TransactionId == "" {
		return errors.New("transactionId is required")
	}

	return c.server.SetTransactionCategory(r.Context(), request.TransactionId, request.CategoryId)
}
