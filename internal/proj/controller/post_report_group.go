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

func (c *reportGroupController) process(r *http.Request) (data.ReportGroupResponse, error) {
	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return data.ReportGroupResponse{}, fmt.Errorf("parsing group report request: %w", err)
	}

	if err := request.Validate(); err != nil {
		return data.ReportGroupResponse{}, err
	}

	// End is used only when it is sent; otherwise the report has no upper bound
	var endMonth, endYear int
	if request.End != nil {
		endMonth, endYear = request.End.Month, request.End.Year
	}

	from, to, err := data.GetDateRange(request.Start.Month, request.Start.Year, endMonth, endYear)
	if err != nil {
		return data.ReportGroupResponse{}, err
	}

	items, err := c.server.GetGroupReport(r.Context(), *from, to, request.ExcludeNone)
	if err != nil {
		return data.ReportGroupResponse{}, err
	}

	return data.ReportGroupResponse{Items: items}, nil
}
