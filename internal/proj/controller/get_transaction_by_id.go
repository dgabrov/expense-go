package controller

import (
	"errors"
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
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, errors.New("id query parameter is required")
	}

	return c.server.GetTransactionById(r.Context(), id)
}
