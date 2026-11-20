package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type ReportCategory interface {
	http.Handler
}

type reportCategoryController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewReportCategory(srv *server.Server, cfg data.ConfigData) ReportCategory {
	return &reportCategoryController{server: srv, cfg: cfg}
}

func (c *reportCategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *reportCategoryController) process(r *http.Request) (string, error) {
	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return "", fmt.Errorf("parsing category report request: %w", err)
	}

	return "hello world", nil
}
