package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type AddTransaction interface {
	http.Handler
}

type addTransactionController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewAddTransaction(srv *server.Server, cfg data.ConfigData) AddTransaction {
	return &addTransactionController{server: srv, cfg: cfg}
}

func (c *addTransactionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *addTransactionController) process(r *http.Request) error {
	var transaction data.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return fmt.Errorf("parsing add transaction request: %w", err)
	}

	return c.server.CreateTransaction(r.Context(), transaction)
}
