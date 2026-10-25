package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
)

type Root interface {
	http.Handler
}

type rootController struct {
	db  *sql.DB
	cfg data.ConfigData
}

func NewRoot(db *sql.DB, cfg data.ConfigData) Root {
	return &rootController{db: db, cfg: cfg}
}

func (c *rootController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "success")
}
