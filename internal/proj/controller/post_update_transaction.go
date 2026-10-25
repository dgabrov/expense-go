package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type UpdateTransaction interface {
	http.Handler
}

type updateTransactionController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewUpdateTransaction(srv *server.Server, cfg data.ConfigData) UpdateTransaction {
	return &updateTransactionController{server: srv, cfg: cfg}
}

func (c *updateTransactionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *updateTransactionController) process(r *http.Request) error {
	var transaction data.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return fmt.Errorf("parsing update transaction request: %w", err)
	}

	return c.server.UpdateTransaction(r.Context(), transaction)
}
