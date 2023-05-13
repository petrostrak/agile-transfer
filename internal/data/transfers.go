package data

import (
	"database/sql"
	"errors"
)

type Transfer struct {
	ID              int64   `json:"id"`
	SourceAccountID int64   `json:"source_account_id"`
	TargetAccountID int64   `json:"target_account_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

type TransferModel struct {
	DB *sql.DB
}

func (t TransferModel) Insert(tx *Transfer) error {
	query := `
		INSERT INTO transfers (source_account_id, target_account_id, amount, currency)
		VALUES ($1, $2, $3)
		RETURNING id, source_account_id, target_account_id, amount, currency`

	args := []any{tx.SourceAccountID, tx.TargetAccountID, tx.Amount, tx.Currency}
	return t.DB.QueryRow(query, args...).Scan(
		&tx.ID,
		&tx.SourceAccountID,
		&tx.TargetAccountID,
		&tx.Amount,
		&tx.Currency,
	)
}

func (t TransferModel) Get(id int64) (*Transfer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, source_account_id, target_account_id, amount, currency
		FROM transfers
		WHERE id = $1`

	var tx Transfer
	err := t.DB.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.SourceAccountID,
		&tx.TargetAccountID,
		&tx.Amount,
		&tx.Currency,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &tx, nil
}
