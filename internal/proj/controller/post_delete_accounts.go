package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type DeleteAccounts interface {
	http.Handler
}

type deleteAccountsController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewDeleteAccounts(srv *server.Server, cfg data.ConfigData) DeleteAccounts {
	return &deleteAccountsController{server: srv, cfg: cfg}
}

func (c *deleteAccountsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *deleteAccountsController) process(r *http.Request) error {
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		return fmt.Errorf("parsing delete accounts request: %w", err)
	}

	// TODO
	return nil
}
