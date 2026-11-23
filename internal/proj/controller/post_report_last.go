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

func (c *reportLastController) process(r *http.Request) ([]data.ReportLastResponse, error) {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return nil, err
	}

	var request data.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, fmt.Errorf("parsing last report request: %w", err)
	}

	if err := request.ValidateStart(); err != nil {
		return nil, err
	}

	// the report starts on the first day of the start month; End and ExcludeNone are ignored
	from, _, err := data.GetDateRange(request.Start.Month, request.Start.Year, 0, 0)
	if err != nil {
		return nil, err
	}

	return c.server.GetLastTransactionDates(r.Context(), *from)
}
