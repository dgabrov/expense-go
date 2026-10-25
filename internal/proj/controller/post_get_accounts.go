package controller

import (
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetAccounts interface {
	http.Handler
}

type getAccountsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetAccounts(srv *server.Server, cfg data.ConfigData) GetAccounts {
	return &getAccountsController{server: srv, cfg: cfg}
}

func (c *getAccountsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getAccountsController) process(r *http.Request) ([]data.Account, error) {
	// TODO
	return nil, nil
}
