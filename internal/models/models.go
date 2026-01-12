package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID      int             `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

type Transaction struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Processed bool            `json:"processed,omitempty"`
}
