package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type ReportList interface {
	http.Handler
}

type reportListController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewReportList(srv *server.Server, cfg data.ConfigData) ReportList {
	return &reportListController{server: srv, cfg: cfg}
}

func (c *reportListController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *reportListController) process(r *http.Request) (data.ReportListResponse, error) {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return data.ReportListResponse{}, err
	}

	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return data.ReportListResponse{}, fmt.Errorf("parsing list report request: %w", err)
	}

	if err := request.Validate(); err != nil {
		return data.ReportListResponse{}, err
	}

	// End is used only when it is sent; otherwise the report has no upper bound
	var endMonth, endYear int
	if request.End != nil {
		endMonth, endYear = request.End.Month, request.End.Year
	}

	from, to, err := data.GetDateRange(request.Start.Month, request.Start.Year, endMonth, endYear)
	if err != nil {
		return data.ReportListResponse{}, err
	}

	items, err := c.server.ListTransactions(r.Context(), *from, to, request.ExcludeNone)
	if err != nil {
		return data.ReportListResponse{}, err
	}

	return data.ReportListResponse{Items: items}, nil
}
