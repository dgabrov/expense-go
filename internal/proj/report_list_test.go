package proj

import (
	"net/http"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
)

const reportListPath = "/expbck/report/list"

func TestReportListStartWithoutEnd(t *testing.T) {
	handler := setupReportHandler(t)

	rec, response := postReport[data.ReportListResponse](t, handler, reportListPath,
		`{"Start":{"Month":11,"Year":2012}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}

	if response.Items == nil {
		t.Fatalf("items is null; expected an array")
	}

	// no End means no upper bound: every transaction must be on or after the start month
	start := monthKey(2012, 11)
	for _, item := range response.Items {
		if monthKey(item.Date.Year, int(item.Date.Month)) < start {
			t.Errorf("transaction %s dated %s is before start 2012-11", item.Id, item.Date)
		}
	}
}

func TestReportListStartAndEnd(t *testing.T) {
	handler := setupReportHandler(t)

	rec, response := postReport[data.ReportListResponse](t, handler, reportListPath,
		`{"Start":{"Month":11,"Year":2012},"End":{"Month":12,"Year":2012}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}

	// End is inclusive at month level: transactions must fall in Nov 2012 through Dec 2012
	start, end := monthKey(2012, 11), monthKey(2012, 12)
	for _, item := range response.Items {
		key := monthKey(item.Date.Year, int(item.Date.Month))
		if key < start || key > end {
			t.Errorf("transaction %s dated %s is outside 2012-11..2012-12", item.Id, item.Date)
		}
	}
}

func TestReportListEndBeforeStart(t *testing.T) {
	handler := setupReportHandler(t)

	assertEndBeforeStartRejected(t, handler, reportListPath)
}
