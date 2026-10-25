package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetTransactionById interface {
	http.Handler
}

type getTransactionByIdController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetTransactionById(srv *server.Server, cfg data.ConfigData) GetTransactionById {
	return &getTransactionByIdController{server: srv, cfg: cfg}
}

func (c *getTransactionByIdController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getTransactionByIdController) process(r *http.Request) (*data.Transaction, error) {
	var id string
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		return nil, fmt.Errorf("parsing get transaction request: %w", err)
	}

	// TODO
	return nil, nil
}
