package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type SearchTransactions interface {
	http.Handler
}

type searchTransactionsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewSearchTransactions(srv *server.Server, cfg data.ConfigData) SearchTransactions {
	return &searchTransactionsController{server: srv, cfg: cfg}
}

func (c *searchTransactionsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *searchTransactionsController) process(r *http.Request) ([]data.Transaction, error) {
	var search data.SearchData
	if err := json.NewDecoder(r.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("parsing search transactions request: %w", err)
	}

	// TODO
	return nil, nil
}
