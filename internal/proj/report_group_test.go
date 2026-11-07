package proj

import (
	"net/http"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
)

const reportGroupPath = "/expbck/report/group"

func TestReportGroupStartWithoutEnd(t *testing.T) {
	handler := setupReportHandler(t)

	rec, response := postReport[data.ReportGroupResponse](t, handler, reportGroupPath,
		`{"Start":{"Month":11,"Year":2012}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}

	if response.Items == nil {
		t.Fatalf("items is null; expected an array")
	}

	// no End means no upper bound: every item must be on or after the start month
	start := monthKey(2012, 11)
	for _, item := range response.Items {
		if monthKey(item.Year, item.Month) < start {
			t.Errorf("item %d-%02d is before start 2012-11", item.Year, item.Month)
		}
	}
}

func TestReportGroupStartAndEnd(t *testing.T) {
	handler := setupReportHandler(t)

	rec, response := postReport[data.ReportGroupResponse](t, handler, reportGroupPath,
		`{"Start":{"Month":11,"Year":2012},"End":{"Month":12,"Year":2012}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}

	// End is inclusive at month level: items must fall in Nov 2012 through Dec 2012
	start, end := monthKey(2012, 11), monthKey(2012, 12)
	for _, item := range response.Items {
		key := monthKey(item.Year, item.Month)
		if key < start || key > end {
			t.Errorf("item %d-%02d is outside 2012-11..2012-12", item.Year, item.Month)
		}
	}
}

func TestReportGroupEndBeforeStart(t *testing.T) {
	handler := setupReportHandler(t)

	assertEndBeforeStartRejected(t, handler, reportGroupPath)
}
