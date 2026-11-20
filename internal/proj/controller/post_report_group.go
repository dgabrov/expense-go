package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type ReportGroup interface {
	http.Handler
}

type reportGroupController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewReportGroup(srv *server.Server, cfg data.ConfigData) ReportGroup {
	return &reportGroupController{server: srv, cfg: cfg}
}

func (c *reportGroupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *reportGroupController) process(r *http.Request) (string, error) {
	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return "", fmt.Errorf("parsing group report request: %w", err)
	}

	return "hello world", nil
}
