package proj

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
)

// setupReportHandler runs the real router against the database named in CONFIG_FILE.
// The test is skipped when the config or the database is not available.
func setupReportHandler(t *testing.T) http.Handler {
	t.Helper()

	if os.Getenv(configEnvVar) == "" {
		t.Skipf("%s not set; skipping HTTP test", configEnvVar)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	db, err := openDB(cfg.Database)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return newHandler(db, cfg)
}

// postReport posts body to path and, on 200, decodes the response into T.
func postReport[T any](t *testing.T, handler http.Handler, path, body string) (*httptest.ResponseRecorder, T) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

// assertEndBeforeStartRejected checks a 400 with the validation message for the given path.
func assertEndBeforeStartRejected(t *testing.T, handler http.Handler, path string) {
	t.Helper()

	rec, _ := postReport[data.ErrorResponse](t, handler, path,
		`{"Start":{"Month":5,"Year":2013},"End":{"Month":1,"Year":2013}}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}

	var errResponse data.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResponse); err != nil {
		t.Fatalf("decoding error response %q: %v", rec.Body.String(), err)
	}
	if errResponse.Message != "start date is after end date" {
		t.Errorf("message = %q, want %q", errResponse.Message, "start date is after end date")
	}
}

// monthKey orders year/month pairs so a range check is a simple comparison.
func monthKey(year, month int) int {
	return year*100 + month
}
