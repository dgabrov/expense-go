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
