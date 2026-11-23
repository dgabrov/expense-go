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

func (c *reportCategoryController) process(r *http.Request) (data.ReportCategoryResponse, error) {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return data.ReportCategoryResponse{}, err
	}

	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return data.ReportCategoryResponse{}, fmt.Errorf("parsing category report request: %w", err)
	}

	if err := request.Validate(); err != nil {
		return data.ReportCategoryResponse{}, err
	}

	// End is used only when it is sent; otherwise the report has no upper bound
	var endMonth, endYear int
	if request.End != nil {
		endMonth, endYear = request.End.Month, request.End.Year
	}

	from, to, err := data.GetDateRange(request.Start.Month, request.Start.Year, endMonth, endYear)
	if err != nil {
		return data.ReportCategoryResponse{}, err
	}

	items, err := c.server.GetCategoryReport(r.Context(), *from, to, request.ExcludeNone)
	if err != nil {
		return data.ReportCategoryResponse{}, err
	}

	return data.ReportCategoryResponse{Items: items}, nil
}
