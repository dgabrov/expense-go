package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
)

type Server struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServer(db *sql.DB, config *data.ConfigData) *Server {
	return &Server{db, config}
}

func (s *Server) CreateSession(ctx context.Context, session data.Session) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO session (session_id, token, user_id, login, name, expiry_dt) VALUES (?, ?, ?, ?, ?, ?)`,
		session.SessionId, session.Token, session.UserId, session.Login, session.Name, session.ExpiryDt)
	if err != nil {
		return fmt.Errorf("inserting session: %w", err)
	}

	return nil
}

func (s *Server) MarkSessionExpired(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE session SET expired_ind = 'Y' WHERE token = ?`, token)
	if err != nil {
		return fmt.Errorf("marking session expired: %w", err)
	}

	return nil
}

func (s *Server) ExtendSession(ctx context.Context, token string, expiryDt time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE session SET expiry_dt = ? WHERE token = ?`, expiryDt, token)
	if err != nil {
		return fmt.Errorf("extending session: %w", err)
	}

	return nil
}

func (s *Server) CreateAccount(ctx context.Context, account data.Account) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO account (account_id, account_name, cash) VALUES (?, ?, ?)`,
		account.Id, account.Name, boolToInd(account.Cash))
	if err != nil {
		return fmt.Errorf("inserting account: %w", err)
	}

	return nil
}

func (s *Server) CreateGroup(ctx context.Context, group data.Group) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO expense_group (group_id, group_name) VALUES (?, ?)`,
		group.Id, group.Name)
	if err != nil {
		return fmt.Errorf("inserting group: %w", err)
	}

	return nil
}

func (s *Server) CreateCategory(ctx context.Context, category data.Category) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO expense_category (category_id, group_id, category_name) VALUES (?, ?, ?)`,
		category.Id, category.GroupId, category.Name)
	if err != nil {
		return fmt.Errorf("inserting category: %w", err)
	}

	return nil
}

func (s *Server) UpdateAccount(ctx context.Context, account data.Account) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE account SET account_name = ?, cash = ? WHERE account_id = ?`,
		account.Name, boolToInd(account.Cash), account.Id)
	if err != nil {
		return fmt.Errorf("updating account: %w", err)
	}

	return checkUpdated(result, "account", account.Id)
}

func (s *Server) UpdateGroup(ctx context.Context, group data.Group) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE expense_group SET group_name = ? WHERE group_id = ?`,
		group.Name, group.Id)
	if err != nil {
		return fmt.Errorf("updating group: %w", err)
	}

	return checkUpdated(result, "group", group.Id)
}

func (s *Server) UpdateCategory(ctx context.Context, category data.Category) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE expense_category SET group_id = ?, category_name = ? WHERE category_id = ?`,
		category.GroupId, category.Name, category.Id)
	if err != nil {
		return fmt.Errorf("updating category: %w", err)
	}

	return checkUpdated(result, "category", category.Id)
}

func (s *Server) UpdateTransaction(ctx context.Context, transaction data.Transaction) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE `transaction` SET name = ?, transaction_date = ?, amount = ?, account_id = ?, category_id = ? WHERE transaction_id = ?",
		transaction.Name, transaction.Date, transaction.Amount, transaction.AccountId, transaction.CategoryId, transaction.Id)
	if err != nil {
		return fmt.Errorf("updating transaction: %w", err)
	}

	return checkUpdated(result, "transaction", transaction.Id)
}

// checkUpdated returns an error when an UPDATE matched no row for the given id.
func checkUpdated(result sql.Result, entity string, id string) error {
	matched, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reading updated rows for %s: %w", entity, err)
	}

	if matched == 0 {
		return fmt.Errorf("%s %s not found", entity, id)
	}

	return nil
}

// boolToInd maps a bool to the Y/N convention used by *ind columns.
func boolToInd(value bool) string {
	if value {
		return "Y"
	}

	return "N"
}

func (s *Server) GetAccounts(ctx context.Context) ([]data.Account, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT account_id, account_name, cash FROM account`)
	if err != nil {
		return nil, fmt.Errorf("querying accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]data.Account, 0)
	for rows.Next() {
		var account data.Account
		var cash string
		if err := rows.Scan(&account.Id, &account.Name, &cash); err != nil {
			return nil, fmt.Errorf("scanning account: %w", err)
		}

		account.Cash = indToBool(cash)
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading accounts: %w", err)
	}

	return accounts, nil
}

func (s *Server) GetGroups(ctx context.Context) ([]data.Group, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT group_id, group_name FROM expense_group`)
	if err != nil {
		return nil, fmt.Errorf("querying groups: %w", err)
	}
	defer rows.Close()

	groups := make([]data.Group, 0)
	for rows.Next() {
		var group data.Group
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return nil, fmt.Errorf("scanning group: %w", err)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading groups: %w", err)
	}

	return groups, nil
}

func (s *Server) GetCategories(ctx context.Context, groupId string) ([]data.Category, error) {
	return s.queryCategories(ctx,
		`SELECT category_id, group_id, category_name FROM expense_category WHERE group_id = ?`, groupId)
}

func (s *Server) GetCompleteGroups(ctx context.Context) ([]data.CompleteGroup, error) {
	groups, err := s.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := s.queryCategories(ctx, `SELECT category_id, group_id, category_name FROM expense_category`)
	if err != nil {
		return nil, err
	}

	result := make([]data.CompleteGroup, 0, len(groups))
	positions := make(map[string]int, len(groups))
	for _, group := range groups {
		positions[group.Id] = len(result)
		result = append(result, data.CompleteGroup{
			Id:         group.Id,
			Name:       group.Name,
			Categories: make([]data.Category, 0),
		})
	}

	for _, category := range categories {
		if position, ok := positions[category.GroupId]; ok {
			result[position].Categories = append(result[position].Categories, category)
		}
	}

	return result, nil
}

// GetTransactionById returns nil, nil when no transaction has the given id.
func (s *Server) GetTransactionById(ctx context.Context, id string) (*data.Transaction, error) {
	var transaction data.Transaction
	var categoryId sql.NullString
	err := s.db.QueryRowContext(ctx,
		"SELECT transaction_id, name, transaction_date, amount, account_id, category_id FROM `transaction` WHERE transaction_id = ?",
		id).Scan(&transaction.Id, &transaction.Name, &transaction.Date, &transaction.Amount, &transaction.AccountId, &categoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying transaction by id: %w", err)
	}

	if categoryId.Valid {
		transaction.CategoryId = &categoryId.String
	}

	return &transaction, nil
}

func (s *Server) queryCategories(ctx context.Context, query string, args ...any) ([]data.Category, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying categories: %w", err)
	}
	defer rows.Close()

	categories := make([]data.Category, 0)
	for rows.Next() {
		var category data.Category
		if err := rows.Scan(&category.Id, &category.GroupId, &category.Name); err != nil {
			return nil, fmt.Errorf("scanning category: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading categories: %w", err)
	}

	return categories, nil
}

// indToBool maps a Y/N *ind column value to a bool; anything other than Y is false.
func indToBool(value string) bool {
	return strings.EqualFold(value, "Y")
}

func (s *Server) CreateTransaction(ctx context.Context, transaction data.Transaction) error {
	return insertTransaction(ctx, s.db, transaction)
}

// CreateTransactions inserts all transactions in one database transaction:
// either every row is stored or none is.
func (s *Server) CreateTransactions(ctx context.Context, transactions []data.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction insert: %w", err)
	}
	defer tx.Rollback()

	for _, transaction := range transactions {
		if err := insertTransaction(ctx, tx, transaction); err != nil {
			return fmt.Errorf("transaction %s: %w", transaction.Id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction insert: %w", err)
	}

	return nil
}

// execer is satisfied by both *sql.DB and *sql.Tx.
type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func insertTransaction(ctx context.Context, e execer, transaction data.Transaction) error {
	_, err := e.ExecContext(ctx,
		"INSERT INTO `transaction` (transaction_id, name, transaction_date, amount, account_id, category_id) VALUES (?, ?, ?, ?, ?, ?)",
		transaction.Id, transaction.Name, transaction.Date, transaction.Amount, transaction.AccountId, transaction.CategoryId)
	if err != nil {
		return fmt.Errorf("inserting transaction: %w", err)
	}

	return nil
}

func (s *Server) DeleteAccounts(ctx context.Context, ids []string) error {
	return s.deleteByIds(ctx, "account", "account_id", "accounts", ids)
}

func (s *Server) DeleteGroups(ctx context.Context, ids []string) error {
	return s.deleteByIds(ctx, "expense_group", "group_id", "groups", ids)
}

func (s *Server) DeleteCategories(ctx context.Context, ids []string) error {
	return s.deleteByIds(ctx, "expense_category", "category_id", "categories", ids)
}

func (s *Server) DeleteTransactions(ctx context.Context, ids []string) error {
	return s.deleteByIds(ctx, "`transaction`", "transaction_id", "transactions", ids)
}

// deleteByIds deletes the rows whose column matches any of ids. It fails and
// deletes nothing when some ids do not exist. plural names the entities in errors.
func (s *Server) deleteByIds(ctx context.Context, table, column, plural string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			args = append(args, id)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting delete of %s: %w", plural, err)
	}
	defer tx.Rollback()

	placeholders := strings.TrimSuffix(strings.Repeat("?, ", len(args)), ", ")
	result, err := tx.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)", table, column, placeholders), args...)
	if err != nil {
		return fmt.Errorf("deleting %s: %w", plural, err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reading deleted %s: %w", plural, err)
	}

	if deleted != int64(len(args)) {
		return fmt.Errorf("only %d of %d %s exist; nothing was deleted", deleted, len(args), plural)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing delete of %s: %w", plural, err)
	}

	return nil
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// SearchTransactions applies only the filters that are set, newest first.
func (s *Server) SearchTransactions(ctx context.Context, search data.SearchData) ([]data.Transaction, error) {
	conditions := make([]string, 0)
	args := make([]any, 0)

	if search.Search != nil && *search.Search != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+likeEscaper.Replace(*search.Search)+"%")
	}
	if search.AccountFilter != nil && *search.AccountFilter != "" {
		conditions = append(conditions, "account_id = ?")
		args = append(args, *search.AccountFilter)
	}
	if search.DateFrom != nil {
		conditions = append(conditions, "transaction_date >= ?")
		args = append(args, *search.DateFrom)
	}
	if search.DateTo != nil {
		conditions = append(conditions, "transaction_date <= ?")
		args = append(args, *search.DateTo)
	}
	if search.AmountFrom != nil {
		conditions = append(conditions, "amount >= ?")
		args = append(args, *search.AmountFrom)
	}
	if search.AmountTo != nil {
		conditions = append(conditions, "amount <= ?")
		args = append(args, *search.AmountTo)
	}

	query := "SELECT transaction_id, name, transaction_date, amount, account_id, category_id FROM `transaction`"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY transaction_date DESC, name"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]data.Transaction, 0)
	for rows.Next() {
		var transaction data.Transaction
		var categoryId sql.NullString
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Date, &transaction.Amount, &transaction.AccountId, &categoryId); err != nil {
			return nil, fmt.Errorf("scanning transaction: %w", err)
		}

		if categoryId.Valid {
			transaction.CategoryId = &categoryId.String
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading transactions: %w", err)
	}

	return transactions, nil
}

// GetLastTransactionDates returns, for each account with a transaction on or after
// from, the account id and its latest transaction date. Ordered by account id.
func (s *Server) GetLastTransactionDates(ctx context.Context, from data.Date) ([]data.ReportLastResponse, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT account_id, MAX(transaction_date) FROM `transaction` WHERE transaction_date >= ? GROUP BY account_id ORDER BY account_id",
		from)
	if err != nil {
		return nil, fmt.Errorf("querying last transaction dates: %w", err)
	}
	defer rows.Close()

	result := make([]data.ReportLastResponse, 0)
	for rows.Next() {
		var last data.ReportLastResponse
		if err := rows.Scan(&last.AccountID, &last.LastDate); err != nil {
			return nil, fmt.Errorf("scanning last transaction date: %w", err)
		}

		result = append(result, last)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading last transaction dates: %w", err)
	}

	return result, nil
}

// GetGroupReport sums transaction amounts per year, month and group, from (inclusive)
// and, when to is not nil, before to (exclusive). Transactions without a category or
// group are reported with an empty group id. When excludeNone is set, the group named
// "none" is left out; unassigned transactions are kept.
func (s *Server) GetGroupReport(ctx context.Context, from data.Date, to *data.Date, excludeNone bool) ([]data.ReportGroupResponseItem, error) {
	conditions := []string{"t.transaction_date >= ?"}
	args := []any{from}

	if to != nil {
		conditions = append(conditions, "t.transaction_date < ?")
		args = append(args, *to)
	}
	// LEFT JOINs keep unassigned transactions, whose group name is NULL; the NULL check
	// is needed because NULL <> 'none' is not true. The collation makes the match case-insensitive.
	if excludeNone {
		conditions = append(conditions, "(g.group_name IS NULL OR g.group_name <> 'none')")
	}

	query := "SELECT YEAR(t.transaction_date) AS report_year, MONTH(t.transaction_date) AS report_month, " +
		"COALESCE(g.group_id, '') AS report_group, SUM(t.amount) AS report_amount " +
		"FROM `transaction` t " +
		"LEFT JOIN expense_category c ON c.category_id = t.category_id " +
		"LEFT JOIN expense_group g ON g.group_id = c.group_id " +
		"WHERE " + strings.Join(conditions, " AND ") + " " +
		"GROUP BY YEAR(t.transaction_date), MONTH(t.transaction_date), COALESCE(g.group_id, '') " +
		"ORDER BY report_year, report_month, report_group"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying group report: %w", err)
	}
	defer rows.Close()

	items := make([]data.ReportGroupResponseItem, 0)
	for rows.Next() {
		var item data.ReportGroupResponseItem
		if err := rows.Scan(&item.Year, &item.Month, &item.GroupID, &item.Amount); err != nil {
			return nil, fmt.Errorf("scanning group report item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading group report: %w", err)
	}

	return items, nil
}

// GetCategoryReport sums transaction amounts per year, month and category, from (inclusive)
// and, when to is not nil, before to (exclusive). Transactions without a category are
// reported with an empty category id. When excludeNone is set, transactions whose category
// belongs to the group named "none" are left out; unassigned transactions are kept.
func (s *Server) GetCategoryReport(ctx context.Context, from data.Date, to *data.Date, excludeNone bool) ([]data.ReportCategoryResponseItem, error) {
	conditions := []string{"t.transaction_date >= ?"}
	args := []any{from}

	if to != nil {
		conditions = append(conditions, "t.transaction_date < ?")
		args = append(args, *to)
	}
	// same NULL-safe group name check as GetGroupReport
	if excludeNone {
		conditions = append(conditions, "(g.group_name IS NULL OR g.group_name <> 'none')")
	}

	query := "SELECT YEAR(t.transaction_date) AS report_year, MONTH(t.transaction_date) AS report_month, " +
		"COALESCE(t.category_id, '') AS report_category, SUM(t.amount) AS report_amount " +
		"FROM `transaction` t " +
		"LEFT JOIN expense_category c ON c.category_id = t.category_id " +
		"LEFT JOIN expense_group g ON g.group_id = c.group_id " +
		"WHERE " + strings.Join(conditions, " AND ") + " " +
		"GROUP BY YEAR(t.transaction_date), MONTH(t.transaction_date), COALESCE(t.category_id, '') " +
		"ORDER BY report_year, report_month, report_category"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying category report: %w", err)
	}
	defer rows.Close()

	items := make([]data.ReportCategoryResponseItem, 0)
	for rows.Next() {
		var item data.ReportCategoryResponseItem
		if err := rows.Scan(&item.Year, &item.Month, &item.CategoryID, &item.Amount); err != nil {
			return nil, fmt.Errorf("scanning category report item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading category report: %w", err)
	}

	return items, nil
}

// ListTransactions returns the transactions dated from (inclusive) and, when to is not nil,
// before to (exclusive), newest first. When excludeNone is set, transactions whose category
// belongs to the group named "none" are left out; unassigned transactions are kept.
func (s *Server) ListTransactions(ctx context.Context, from data.Date, to *data.Date, excludeNone bool) ([]data.Transaction, error) {
	conditions := []string{"t.transaction_date >= ?"}
	args := []any{from}

	if to != nil {
		conditions = append(conditions, "t.transaction_date < ?")
		args = append(args, *to)
	}
	// same NULL-safe group name check as GetGroupReport
	if excludeNone {
		conditions = append(conditions, "(g.group_name IS NULL OR g.group_name <> 'none')")
	}

	query := "SELECT t.transaction_id, t.name, t.transaction_date, t.amount, t.account_id, t.category_id " +
		"FROM `transaction` t " +
		"LEFT JOIN expense_category c ON c.category_id = t.category_id " +
		"LEFT JOIN expense_group g ON g.group_id = c.group_id " +
		"WHERE " + strings.Join(conditions, " AND ") + " " +
		"ORDER BY t.transaction_date DESC, t.name"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying list report: %w", err)
	}
	defer rows.Close()

	transactions := make([]data.Transaction, 0)
	for rows.Next() {
		var transaction data.Transaction
		var categoryId sql.NullString
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Date, &transaction.Amount, &transaction.AccountId, &categoryId); err != nil {
			return nil, fmt.Errorf("scanning list report transaction: %w", err)
		}

		if categoryId.Valid {
			transaction.CategoryId = &categoryId.String
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading list report: %w", err)
	}

	return transactions, nil
}

// GetMonthGroupTransactions returns the transactions dated from (inclusive) to to (exclusive).
// When groupId is empty, only transactions without a category are returned; otherwise only
// transactions whose category belongs to the group groupId.
func (s *Server) GetMonthGroupTransactions(ctx context.Context, from data.Date, to data.Date, groupId string) ([]data.Transaction, error) {
	conditions := []string{"t.transaction_date >= ?", "t.transaction_date < ?"}
	args := []any{from, to}

	if groupId == "" {
		conditions = append(conditions, "t.category_id IS NULL")
	} else {
		conditions = append(conditions, "t.category_id IN (SELECT c.category_id FROM expense_category c WHERE c.group_id = ?)")
		args = append(args, groupId)
	}

	query := "SELECT t.transaction_id, t.name, t.transaction_date, t.amount, t.account_id, t.category_id " +
		"FROM `transaction` t " +
		"WHERE " + strings.Join(conditions, " AND ") + " " +
		"ORDER BY t.transaction_date DESC, t.name"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying month group transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]data.Transaction, 0)
	for rows.Next() {
		var transaction data.Transaction
		var categoryId sql.NullString
		if err := rows.Scan(&transaction.Id, &transaction.Name, &transaction.Date, &transaction.Amount, &transaction.AccountId, &categoryId); err != nil {
			return nil, fmt.Errorf("scanning month group transaction: %w", err)
		}

		if categoryId.Valid {
			transaction.CategoryId = &categoryId.String
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading month group transactions: %w", err)
	}

	return transactions, nil
}

// GetSessionByToken returns nil, nil when no session has the given token.
func (s *Server) GetSessionByToken(ctx context.Context, token string) (*data.Session, error) {
	var session data.Session
	err := s.db.QueryRowContext(ctx,
		`SELECT session_id, token, user_id, login, name, expiry_dt, expired_ind = 'Y' FROM session WHERE token = ?`,
		token).Scan(&session.SessionId, &session.Token, &session.UserId, &session.Login, &session.Name, &session.ExpiryDt, &session.Expired)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying session by token: %w", err)
	}

	return &session, nil
}
