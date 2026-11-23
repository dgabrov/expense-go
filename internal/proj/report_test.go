package proj

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
	"github.com/google/uuid"
)

// sessionCookieName must match the cookie name used by the controllers.
const sessionCookieName = "expenseCookie12"

// setupReportHandler runs the real router against the database named in CONFIG_FILE and
// returns it with a logged-in session cookie. The session is expired at cleanup.
// The test is skipped when the config or the database is not available.
func setupReportHandler(t *testing.T) (http.Handler, *http.Cookie) {
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

	srv := server.NewServer(db, &cfg)

	token := uuid.NewString()
	session := data.Session{
		SessionId: uuid.NewString(),
		Token:     token,
		UserId:    "report-test",
		Login:     "report-test",
		Name:      "report test",
		ExpiryDt:  time.Now().UTC().Add(time.Duration(cfg.TokenTTLSeconds) * time.Second),
	}
	if err := srv.CreateSession(context.Background(), session); err != nil {
		t.Fatalf("creating test session: %v", err)
	}
	t.Cleanup(func() {
		if err := srv.MarkSessionExpired(context.Background(), token); err != nil {
			t.Errorf("expiring test session: %v", err)
		}
	})

	return newHandler(db, cfg), &http.Cookie{Name: sessionCookieName, Value: token}
}

// postReport posts body to path with the session cookie and, on 200, decodes the response into T.
func postReport[T any](t *testing.T, handler http.Handler, cookie *http.Cookie, path, body string) (*httptest.ResponseRecorder, T) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

// assertEndBeforeStartRejected checks a 400 with the validation message for the given path.
func assertEndBeforeStartRejected(t *testing.T, handler http.Handler, cookie *http.Cookie, path string) {
	t.Helper()

	rec, _ := postReport[data.ErrorResponse](t, handler, cookie, path,
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

func TestReportRejectsMissingSession(t *testing.T) {
	handler, _ := setupReportHandler(t)

	rec, _ := postReport[data.ErrorResponse](t, handler, nil, reportGroupPath,
		`{"Start":{"Month":11,"Year":2012}}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "not logged in") {
		t.Errorf("body %s does not mention not logged in", rec.Body.String())
	}
}
