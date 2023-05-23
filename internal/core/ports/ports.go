package ports

import (
	"context"

	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Insert(acc *domain.Account) error
	Get(id int64) (*domain.Account, error)
	Update(account *domain.Account) error
	Delete(id int64) error
	GetAll(ctx context.Context) ([]domain.Account, error)
}

type TransferRepository interface {
	Insert(ctx context.Context, tx domain.Transfer) (domain.Transfer, error)
	Get(id int64) (*domain.Transfer, error)
	GetAll() ([]domain.Transfer, error)
	ExecTx(ctx context.Context, fn func() error) error
	TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error)
	AddAccountBalance(ctx context.Context, id int64, amount decimal.Decimal) (domain.Account, error)
	AddMoney(ctx context.Context, sourceAccountID int64, sourceAccountAmount decimal.Decimal, targetAccountID int64, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount domain.Account, err error)
	ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]domain.Account, error)
}
