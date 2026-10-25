package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type AddAccount interface {
	http.Handler
}

type addAccountController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewAddAccount(srv *server.Server, cfg data.ConfigData) AddAccount {
	return &addAccountController{server: srv, cfg: cfg}
}

func (c *addAccountController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *addAccountController) process(r *http.Request) error {
	var account data.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return fmt.Errorf("parsing add account request: %w", err)
	}

	return c.server.CreateAccount(r.Context(), account)
}
