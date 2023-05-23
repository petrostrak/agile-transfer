package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID              int64           `json:"id"`
	SourceAccountID int64           `json:"source_account_id"`
	TargetAccountID int64           `json:"target_account_id"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
}

type transferRepository struct {
	DB *sql.DB
}

func (t *transferRepository) Insert(ctx context.Context, tx Transfer) (Transfer, error) {
	query := `
		INSERT INTO transfers (source_account_id, target_account_id, amount, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, source_account_id, target_account_id, amount, currency`

	args := []any{tx.SourceAccountID, tx.TargetAccountID, tx.Amount, tx.Currency}
	var transfer Transfer
	err := t.DB.QueryRowContext(ctx, query, args...).Scan(
		&transfer.ID,
		&transfer.SourceAccountID,
		&transfer.TargetAccountID,
		&transfer.Amount,
		&transfer.Currency,
	)

	return transfer, err
}

func (t *transferRepository) Get(id int64) (*Transfer, error) {
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

func (t *transferRepository) GetAll() ([]Transfer, error) {
	query := `
			SELECT id, source_account_id, target_account_id, amount, currency
			FROM transfers
			ORDER BY id`

	rows, err := t.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var transfer Transfer
		if err := rows.Scan(
			&transfer.ID,
			&transfer.SourceAccountID,
			&transfer.TargetAccountID,
			&transfer.Amount,
			&transfer.Currency,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (t *transferRepository) ExecTx(ctx context.Context, fn func() error) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			return ctx.Err()
		default:
			return err
		}
	}

	if err = fn(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}