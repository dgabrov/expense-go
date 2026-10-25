package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/dgabrov/expense-go/internal/proj/server"
	"github.com/google/uuid"
)

const (
	authTimeout       = 10 * time.Second
	sessionCookieName = "expenseCookie12"
)

var authClient = &http.Client{Timeout: authTimeout}

type Login interface {
	http.Handler
}

type loginController struct {
	server *server.Server
	cfg    data.ConfigData
}

func NewLogin(srv *server.Server, cfg data.ConfigData) Login {
	return &loginController{server: srv, cfg: cfg}
}

func (c *loginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := c.process(w, r)

	writeJsonResponse(w, payload, err)
}

func (c *loginController) process(w http.ResponseWriter, r *http.Request) (*data.LoggedUser, error) {
	var login data.LoginData
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		return nil, fmt.Errorf("parsing login request: %w", err)
	}

	body, err := json.Marshal(login)
	if err != nil {
		return nil, fmt.Errorf("encoding auth request: %w", err)
	}

	authReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, c.cfg.Auth, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building auth request: %w", err)
	}
	authReq.Header.Set("Content-Type", "application/json")

	authResp, err := authClient.Do(authReq)
	if err != nil {
		return nil, fmt.Errorf("calling auth service: %w", err)
	}
	defer authResp.Body.Close()

	if authResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", authResp.StatusCode)
	}

	var user data.AuthLoginResponse
	if err := json.NewDecoder(authResp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("parsing auth response: %w", err)
	}

	session := data.Session{
		SessionId: uuid.Must(uuid.NewV7()).String(),
		Token:     createRandomToken(),
		UserId:    user.Id,
		Login:     user.Login,
		Name:      user.Name,
		ExpiryDt:  time.Now().UTC().Add(time.Duration(c.cfg.TokenTTLSeconds) * time.Second),
	}

	if err := c.server.CreateSession(r.Context(), session); err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		MaxAge:   c.cfg.TokenTTLSeconds,
		HttpOnly: true,
	})

	return &data.LoggedUser{
		Id:    user.Id,
		Login: user.Login,
		Name:  user.Name,
	}, nil
}
