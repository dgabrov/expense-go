package controller

import (
	"database/sql"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
)

type Login interface {
	PostLogin() http.HandlerFunc
}

type loginController struct {
	db  *sql.DB
	cfg data.ConfigData
}

func NewLogin(db *sql.DB, cfg data.ConfigData) Login {
	return &loginController{db: db, cfg: cfg}
}

func (c *loginController) PostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
