package data

import "time"

type Session struct {
	SessionId string
	Token     string
	UserId    string
	Login     string
	Name      string
	ExpiryDt  time.Time
}
