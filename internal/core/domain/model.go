package domain

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

type TransferTxParams struct {
	SourceAccountID  int64           `json:"source_account_id"`
	TargetAccountID  int64           `json:"target_account_id"`
	SourceBalance    decimal.Decimal `json:"source_balance"`
	AmountToTransfer decimal.Decimal `json:"amount_to_transfer"`
	SourceCurrency   string          `json:"source_currency"`
	TargetCurrency   string          `json:"target_currency"`
}

type TransferTxResult struct {
	Transfer      `json:"transfer"`
	SourceAccount Account `json:"source_account"`
	TargetAccount Account `json:"target_account"`
}
