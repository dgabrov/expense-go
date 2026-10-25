package controller

import (
	"database/sql"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
)

type Login interface {
	http.Handler
}

type loginController struct {
	db  *sql.DB
	cfg data.ConfigData
}

func NewLogin(db *sql.DB, cfg data.ConfigData) Login {
	return &loginController{db: db, cfg: cfg}
}

func (c *loginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *loginController) process(r *http.Request) (*data.LoggedUser, error) {
	return nil, nil
}
