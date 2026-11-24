package proj

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
)

const monthGroupPath = "/expbck/month/group"

// getMonthGroup sends a GET with the given query string and, on 200, decodes the response into T.
func getMonthGroup[T any](t *testing.T, handler http.Handler, cookie *http.Cookie, query string) (*httptest.ResponseRecorder, T) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, monthGroupPath+"?"+query, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var response T
	if rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("decoding response %q: %v", rec.Body.String(), err)
		}
	}

	return rec, response
}

func TestMonthGroupRejectsMissingSession(t *testing.T) {
	handler, _ := setupReportHandler(t)

	rec, _ := getMonthGroup[data.ErrorResponse](t, handler, nil, "month=11&year=2012&groupId=")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "not logged in") {
		t.Errorf("body %s does not mention not logged in", rec.Body.String())
	}
}

func TestMonthGroupRejectsInvalidParameters(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	cases := map[string]string{
		"month missing":  "year=2012&groupId=",
		"month not num":  "month=ab&year=2012&groupId=",
		"month zero":     "month=0&year=2012&groupId=",
		"month thirteen": "month=13&year=2012&groupId=",
		"year missing":   "month=11&groupId=",
		"year not num":   "month=11&year=abc&groupId=",
		"year zero":      "month=11&year=0&groupId=",
	}

	for name, query := range cases {
		t.Run(name, func(t *testing.T) {
			rec, _ := getMonthGroup[data.ErrorResponse](t, handler, cookie, query)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
			}
		})
	}
}

// With an empty groupId every returned transaction has no category, and it is within the month.
func TestMonthGroupEmptyGroupReturnsUncategorized(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	rec, response := getMonthGroup[[]data.Transaction](t, handler, cookie, "month=11&year=2012&groupId=")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}
	if response == nil {
		t.Fatalf("response is null; expected an array")
	}

	for _, transaction := range response {
		if transaction.CategoryId != nil {
			t.Errorf("transaction %s has category %s; expected none", transaction.Id, *transaction.CategoryId)
		}
		if transaction.Date.Year != 2012 || transaction.Date.Month != 11 {
			t.Errorf("transaction %s dated %s is outside 2012-11", transaction.Id, transaction.Date)
		}
	}
}

// A groupId that matches no group yields an empty array, not null.
func TestMonthGroupUnknownGroupReturnsEmptyArray(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	rec, response := getMonthGroup[[]data.Transaction](t, handler, cookie, "month=11&year=2012&groupId=no-such-group")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}
	if response == nil {
		t.Fatalf("response is null; expected an empty array")
	}
	if len(response) != 0 {
		t.Errorf("got %d transactions, want 0", len(response))
	}
}
