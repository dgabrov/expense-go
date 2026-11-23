package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type AddBulk interface {
	http.Handler
}

type addBulkController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewAddBulk(srv *server.Server, cfg data.ConfigData) AddBulk {
	return &addBulkController{server: srv, cfg: cfg}
}

func (c *addBulkController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *addBulkController) process(r *http.Request) error {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return err
	}

	var transactions []data.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactions); err != nil {
		return fmt.Errorf("parsing add bulk request: %w", err)
	}

	return c.server.CreateTransactions(r.Context(), transactions)
}
