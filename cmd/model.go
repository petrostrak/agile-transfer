package cmd

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
}

type Transfer struct {
	ID              int64           `json:"id"`
	SourceAccountID int64           `json:"source_account_id"`
	TargetAccountID int64           `json:"target_account_id"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
}
