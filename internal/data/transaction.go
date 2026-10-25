package data

import "time"

type Transaction struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	Amount     float64   `json:"amount"`
	AccountId  string    `json:"accountId"`
	CategoryId *string   `json:"categoryId"`
}

type SearchData struct {
	Search        *string    `json:"search"`
	AccountFilter *string    `json:"accountFilter"`
	DateFrom      *time.Time `json:"dateFrom"`
	DateTo        *time.Time `json:"dateTo"`
	AmountFrom    *float64   `json:"amountFrom"`
	AmountTo      *float64   `json:"amountTo"`
}
