package proj

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
)

const monthCategoryPath = "/expbck/month/category"

// getMonthCategory sends a GET with the given query string and, on 200, decodes the response into T.
func getMonthCategory[T any](t *testing.T, handler http.Handler, cookie *http.Cookie, query string) (*httptest.ResponseRecorder, T) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, monthCategoryPath+"?"+query, nil)
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

func TestMonthCategoryRejectsMissingSession(t *testing.T) {
	handler, _ := setupReportHandler(t)

	rec, _ := getMonthCategory[data.ErrorResponse](t, handler, nil, "month=11&year=2012&categoryId=")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "not logged in") {
		t.Errorf("body %s does not mention not logged in", rec.Body.String())
	}
}

func TestMonthCategoryRejectsInvalidParameters(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	cases := map[string]string{
		"month missing":  "year=2012&categoryId=",
		"month not num":  "month=ab&year=2012&categoryId=",
		"month zero":     "month=0&year=2012&categoryId=",
		"month thirteen": "month=13&year=2012&categoryId=",
		"year missing":   "month=11&categoryId=",
		"year not num":   "month=11&year=abc&categoryId=",
		"year zero":      "month=11&year=0&categoryId=",
	}

	for name, query := range cases {
		t.Run(name, func(t *testing.T) {
			rec, _ := getMonthCategory[data.ErrorResponse](t, handler, cookie, query)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
			}
		})
	}
}

// With an empty categoryId every returned transaction has no category, and it is within the month.
func TestMonthCategoryEmptyReturnsUncategorized(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	rec, response := getMonthCategory[[]data.Transaction](t, handler, cookie, "month=11&year=2012&categoryId=")
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

// A categoryId that matches no category yields an empty array, not null.
func TestMonthCategoryUnknownCategoryReturnsEmptyArray(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	rec, response := getMonthCategory[[]data.Transaction](t, handler, cookie, "month=11&year=2012&categoryId=no-such-category")
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
