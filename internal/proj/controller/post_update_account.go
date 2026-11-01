package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type UpdateAccount interface {
	http.Handler
}

type updateAccountController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewUpdateAccount(srv *server.Server, cfg data.ConfigData) UpdateAccount {
	return &updateAccountController{server: srv, cfg: cfg}
}

func (c *updateAccountController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.process(r)

	writeJsonResponse(w, nil, err)
}

func (c *updateAccountController) process(r *http.Request) error {
	var account data.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return fmt.Errorf("parsing update account request: %w", err)
	}

	return c.server.UpdateAccount(r.Context(), account)
}
