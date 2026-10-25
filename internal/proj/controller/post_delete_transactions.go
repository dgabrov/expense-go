package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type DeleteTransactions interface {
	http.Handler
}

type deleteTransactionsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewDeleteTransactions(srv *server.Server, cfg data.ConfigData) DeleteTransactions {
	return &deleteTransactionsController{server: srv, cfg: cfg}
}

func (c *deleteTransactionsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *deleteTransactionsController) process(r *http.Request) error {
	var holder data.IdsHolder
	if err := json.NewDecoder(r.Body).Decode(&holder); err != nil {
		return fmt.Errorf("parsing delete transactions request: %w", err)
	}

	return c.server.DeleteTransactions(r.Context(), holder.Ids)
}
