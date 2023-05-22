package main

import (
	"context"

	"github.com/shopspring/decimal"
)

type AccountService interface {
	Insert(acc *Account) error
	Get(id int64) (*Account, error)
	Update(account *Account) error
	Delete(id int64) error
	GetAll(ctx context.Context) ([]Account, error)
	AddAccountBalance(ctx context.Context, id int64, amount decimal.Decimal) (Account, error)
	AddMoney(ctx context.Context, sourceAccountID int64, sourceAccountAmount decimal.Decimal, targetAccountID int64, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount Account, err error)
	ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]Account, error)
}

type TransferService interface {
	Insert(ctx context.Context, tx Transfer) (Transfer, error)
	Get(id int64) (*Transfer, error)
	GetAll() ([]Transfer, error)
	ExecTx(ctx context.Context, fn func() error) error
	TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error)
}
