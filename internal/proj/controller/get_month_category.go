package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
)

type GetMonthCategory interface {
	http.Handler
}

type getMonthCategoryController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewGetMonthCategory(srv *server.Server, cfg data.ConfigData) GetMonthCategory {
	return &getMonthCategoryController{server: srv, cfg: cfg}
}

func (c *getMonthCategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(r)

	writeJsonResponse(w, payload, err)
}

func (c *getMonthCategoryController) process(r *http.Request) ([]data.Transaction, error) {
	if err := checkLogged(r, c.server, c.cfg); err != nil {
		return nil, err
	}

	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		return nil, fmt.Errorf("month query parameter must be a number from 1 to 12: %w", err)
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		return nil, fmt.Errorf("year query parameter must be a number: %w", err)
	}

	// categoryId may be empty: that selects the transactions without a category
	categoryId := query.Get("categoryId")

	// the same start and end month yields the first day of the month and the first day of the next month
	from, to, err := data.GetDateRange(month, year, month, year)
	if err != nil {
		return nil, err
	}
	if to == nil {
		return nil, errors.New("month range has no end")
	}

	return c.server.GetMonthCategoryTransactions(r.Context(), *from, *to, categoryId)
}
