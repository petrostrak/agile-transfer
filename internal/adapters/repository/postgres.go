package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/shopspring/decimal"
)

var (
	ErrRecordNotFound = errors.New("record not Found")
	POSTGRES_DSN      = "postgres://postgres:password@localhost/agile_transfer?sslmode=disable"
)

type PostgresRepository struct {
	*AccountRepository
	*TransferRepository
}

func NewPostgressRepository() *PostgresRepository {
	dsn := POSTGRES_DSN
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil
	}

	return &PostgresRepository{
		&AccountRepository{db},
		&TransferRepository{db},
	}
}

type TransferRepository struct {
	DB *sql.DB
}

func (t *TransferRepository) Insert(ctx context.Context, tx domain.Transfer) (domain.Transfer, error) {
	query := `
		INSERT INTO transfers (source_account_id, target_account_id, amount, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, source_account_id, target_account_id, amount, currency`

	args := []any{tx.SourceAccountID, tx.TargetAccountID, tx.Amount, tx.Currency}
	var transfer domain.Transfer
	err := t.DB.QueryRowContext(ctx, query, args...).Scan(
		&transfer.ID,
		&transfer.SourceAccountID,
		&transfer.TargetAccountID,
		&transfer.Amount,
		&transfer.Currency,
	)

	return transfer, err
}

func (t *TransferRepository) Get(id uuid.UUID) (*domain.Transfer, error) {
	query := `
		SELECT id, source_account_id, target_account_id, amount, currency
		FROM transfers
		WHERE id = $1`

	var tx domain.Transfer
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

func (t *TransferRepository) GetAll() ([]domain.Transfer, error) {
	query := `
			SELECT id, source_account_id, target_account_id, amount, currency
			FROM transfers
			ORDER BY id`

	rows, err := t.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []domain.Transfer
	for rows.Next() {
		var transfer domain.Transfer
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

func (t *TransferRepository) ExecTx(ctx context.Context, fn func() error) error {
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

func (t *TransferRepository) AddAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (domain.Account, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING id, balance, currency, created_at`

	args := []any{amount, id}

	var account domain.Account
	err := t.DB.QueryRowContext(ctx, query, args...).Scan(
		&account.ID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)

	return account, err
}

func (t *TransferRepository) AddMoney(ctx context.Context, sourceAccountID uuid.UUID, sourceAccountAmount decimal.Decimal, targetAccountID uuid.UUID, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount domain.Account, err error) {
	sourceAccount, err = t.AddAccountBalance(ctx, sourceAccountID, sourceAccountAmount)
	if err != nil {
		return
	}

	targetAccount, err = t.AddAccountBalance(ctx, targetAccountID, targetAccountAmount)
	if err != nil {
		return
	}

	return
}

func (t *TransferRepository) TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error) {
	var result domain.TransferTxResult

	err := t.ExecTx(ctx, func() error {
		var err error

		result.SourceAccount, result.TargetAccount, err = t.AddMoney(
			ctx,
			arg.SourceAccountID,
			arg.AmountToTransfer.Neg(),
			arg.TargetAccountID,
			arg.AmountToTransfer,
		)
		if err != nil {
			return err
		}

		trasfer := domain.Transfer{
			SourceAccountID: arg.SourceAccountID,
			TargetAccountID: arg.TargetAccountID,
			Amount:          arg.AmountToTransfer,
			Currency:        arg.TargetCurrency,
		}

		result.Transfer, err = t.Insert(ctx, trasfer)
		if err != nil {
			return err
		}

		return nil
	})

	return &result, err
}

func (t *TransferRepository) ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID uuid.UUID) ([]domain.Account, error) {
	query := `
		SELECT id, balance, currency, created_at 
		FROM accounts
		WHERE id IN ($1, $2)`

	args := []any{sourceAccountID, targetAccountID}

	rows, err := t.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if len(accounts) != 2 {
		return nil, errors.New("one or more of the accounts given does not exist")
	}

	return accounts, nil
}

type AccountRepository struct {
	DB *sql.DB
}

type Account struct {
	ID        uuid.UUID       `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
}

func (a *AccountRepository) Insert(acc *domain.Account) error {
	query := `
		INSERT INTO accounts (balance, currency)
		VALUES ($1, $2)
		RETURNING id, balance, currency, created_at`

	args := []any{acc.Balance, acc.Currency}

	return a.DB.QueryRow(query, args...).Scan(&acc.ID, &acc.Balance, &acc.Currency, &acc.CreatedAt)
}

func (a *AccountRepository) Get(id uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, balance, currency, created_at 
		FROM accounts
		WHERE id = $1`

	var account domain.Account

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

func (a *AccountRepository) Update(account *domain.Account) error {
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

func (a *AccountRepository) Delete(id uuid.UUID) error {
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

func (a *AccountRepository) GetAll(ctx context.Context) ([]domain.Account, error) {
	query := `
		SELECT id, balance, currency, created_at
		FROM accounts
		ORDER BY id`

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
