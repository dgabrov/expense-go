package proj

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
	"github.com/google/uuid"
)

const setCategoryPath = "/expbck/setCategory"

// putSetCategory sends a PUT with body and the optional session cookie and returns the recorder.
func putSetCategory(handler http.Handler, cookie *http.Cookie, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, setCategoryPath, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	return rec
}

// setCategoryFixture creates an account, a group, a category and a transaction in the
// database and removes them at cleanup. The category belongs to the group.
func setCategoryFixture(t *testing.T) (ids struct{ transactionId, categoryId string }, db *sql.DB) {
	t.Helper()

	if os.Getenv(configEnvVar) == "" {
		t.Skipf("%s not set; skipping HTTP test", configEnvVar)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	db, err = openDB(cfg.Database)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	srv := server.NewServer(db, &cfg)
	ctx := context.Background()

	accountId := uuid.NewString()
	groupId := uuid.NewString()
	ids.categoryId = uuid.NewString()
	ids.transactionId = uuid.NewString()

	if err := srv.CreateAccount(ctx, data.Account{Id: accountId, Name: "set category test"}); err != nil {
		t.Fatalf("creating account: %v", err)
	}
	if err := srv.CreateGroup(ctx, data.Group{Id: groupId, Name: "set category test"}); err != nil {
		t.Fatalf("creating group: %v", err)
	}
	if err := srv.CreateCategory(ctx, data.Category{Id: ids.categoryId, GroupId: groupId, Name: "set category test"}); err != nil {
		t.Fatalf("creating category: %v", err)
	}
	if err := srv.CreateTransaction(ctx, data.Transaction{
		Id:        ids.transactionId,
		Name:      "set category test",
		Date:      data.Date{Year: 2012, Month: 11, Day: 1},
		Amount:    1,
		AccountId: accountId,
	}); err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	t.Cleanup(func() {
		// children first, because of the foreign keys
		for _, stmt := range []struct{ query, id string }{
			{"DELETE FROM `transaction` WHERE transaction_id = ?", ids.transactionId},
			{"DELETE FROM expense_category WHERE category_id = ?", ids.categoryId},
			{"DELETE FROM expense_group WHERE group_id = ?", groupId},
			{"DELETE FROM account WHERE account_id = ?", accountId},
		} {
			if _, err := db.Exec(stmt.query, stmt.id); err != nil {
				t.Errorf("cleaning up %q: %v", stmt.query, err)
			}
		}
	})

	return ids, db
}

// transactionCategory reads the category_id column of a transaction; nil means NULL.
func transactionCategory(t *testing.T, db *sql.DB, transactionId string) *string {
	t.Helper()

	var category sql.NullString
	if err := db.QueryRow("SELECT category_id FROM `transaction` WHERE transaction_id = ?", transactionId).Scan(&category); err != nil {
		t.Fatalf("reading category of transaction %s: %v", transactionId, err)
	}
	if !category.Valid {
		return nil
	}
	return &category.String
}

func TestSetCategoryRejectsMissingSession(t *testing.T) {
	handler, _ := setupReportHandler(t)

	rec := putSetCategory(handler, nil, `{"transactionId":"x","categoryId":""}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "not logged in") {
		t.Errorf("body %s does not mention not logged in", rec.Body.String())
	}
}

func TestSetCategoryRejectsMissingTransactionId(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	rec := putSetCategory(handler, cookie, `{"categoryId":"x"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
}

func TestSetCategoryRejectsUnknownTransaction(t *testing.T) {
	handler, cookie := setupReportHandler(t)

	body := `{"transactionId":"` + uuid.NewString() + `","categoryId":""}`
	rec := putSetCategory(handler, cookie, body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "not found") {
		t.Errorf("body %s does not mention not found", rec.Body.String())
	}
}

func TestSetCategoryRejectsUnknownCategory(t *testing.T) {
	handler, cookie := setupReportHandler(t)
	ids, _ := setCategoryFixture(t)

	body := `{"transactionId":"` + ids.transactionId + `","categoryId":"` + uuid.NewString() + `"}`
	rec := putSetCategory(handler, cookie, body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "does not exist") {
		t.Errorf("body %s does not mention a missing category", rec.Body.String())
	}
}

func TestSetCategorySetsAndClears(t *testing.T) {
	handler, cookie := setupReportHandler(t)
	ids, db := setCategoryFixture(t)

	// set the category
	body := `{"transactionId":"` + ids.transactionId + `","categoryId":"` + ids.categoryId + `"}`
	rec := putSetCategory(handler, cookie, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("set: status = %d, body %s", rec.Code, rec.Body.String())
	}
	if got := transactionCategory(t, db, ids.transactionId); got == nil || *got != ids.categoryId {
		t.Fatalf("after set: category = %v, want %s", got, ids.categoryId)
	}

	// an empty categoryId clears it
	body = `{"transactionId":"` + ids.transactionId + `","categoryId":""}`
	rec = putSetCategory(handler, cookie, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("clear: status = %d, body %s", rec.Code, rec.Body.String())
	}
	if got := transactionCategory(t, db, ids.transactionId); got != nil {
		t.Fatalf("after clear: category = %s, want NULL", *got)
	}
}
