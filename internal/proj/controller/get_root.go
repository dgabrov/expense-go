package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
)

type Root interface {
	GetRoot() http.HandlerFunc
}

type rootController struct {
	db  *sql.DB
	cfg data.ConfigData
}

func NewRoot(db *sql.DB, cfg data.ConfigData) Root {
	return &rootController{db: db, cfg: cfg}
}

func (c *rootController) GetRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "success")
	}
}
