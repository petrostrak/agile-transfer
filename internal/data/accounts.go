package data

import (
	"database/sql"
	"errors"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountModel struct {
	DB *sql.DB
}

func (a AccountModel) Insert(acc *Account) error {
	query := `
		INSERT INTO accounts (balance, currency)
		VALUES ($1, $2)
		RETURNING id, balance, currency, created_at`

	args := []any{acc.Balance, acc.Currency}

	return a.DB.QueryRow(query, args...).Scan(&acc.ID, &acc.Balance, &acc.Currency, &acc.CreatedAt)
}

func (a AccountModel) Get(id int64) (*Account, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, balance, currency, created_at 
		FROM accounts
		WHERE id = $1`

	var account Account

	err := a.DB.QueryRow(query, id).Scan(
		&account.ID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &account, nil
}

func (a AccountModel) Update(account *Account) error {
	query := `
		UPDATE accounts
		SET balance = $1
		WHERE id = $2
		RETURNING id, balance, currency, created_at`

	args := []any{account.Balance, account.ID}

	return a.DB.QueryRow(query, args...).Scan(
		&account.ID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
}

func (a AccountModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM accounts
		WHERE id = $1`

	result, err := a.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (a AccountModel) AddAccountBalance(id int64, amount float64) (Account, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING id, balance, currency, created_at`

	args := []any{amount, id}

	var account Account
	err := a.DB.QueryRow(query, args...).Scan(
		&account.ID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)

	return account, err
}

func (a AccountModel) AddMoney(sourceAccountID int64, sourceAccountAmount float64, targetAccountID int64, targetAccountAmount float64) (sourceAccount, targetAccount Account, err error) {
	sourceAccount, err = a.AddAccountBalance(sourceAccountID, sourceAccountAmount)
	if err != nil {
		return
	}

	targetAccount, err = a.AddAccountBalance(targetAccountID, targetAccountAmount)
	if err != nil {
		return
	}

	return
}
