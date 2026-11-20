package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type ReportLast interface {
	http.Handler
}

type reportLastController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewReportLast(srv *server.Server, cfg data.ConfigData) ReportLast {
	return &reportLastController{server: srv, cfg: cfg}
}

func (c *reportLastController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *reportLastController) process(r *http.Request) (string, error) {
	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return "", fmt.Errorf("parsing last report request: %w", err)
	}

	return "hello world", nil
}
