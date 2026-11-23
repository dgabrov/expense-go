package controller

import (
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type Root interface {
	http.Handler
}

type rootController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewRoot(srv *server.Server, cfg data.ConfigData) Root {
	return &rootController{server: srv, cfg: cfg}
}

func (c *rootController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		writeJsonResponse(w, nil, err)
		return
	}

	fmt.Fprint(w, "success")
}
