package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Insert(acc *domain.Account) error
	Get(id uuid.UUID) (*domain.Account, error)
	Update(account *domain.Account) error
	Delete(id uuid.UUID) error
	GetAll(ctx context.Context) ([]domain.Account, error)
}

type TransferRepository interface {
	Insert(ctx context.Context, tx domain.Transfer) (domain.Transfer, error)
	Get(id uuid.UUID) (*domain.Transfer, error)
	GetAll() ([]domain.Transfer, error)
	ExecTx(ctx context.Context, fn func() error) error
	TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error)
	AddAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (domain.Account, error)
	AddMoney(ctx context.Context, sourceAccountID uuid.UUID, sourceAccountAmount decimal.Decimal, targetAccountID uuid.UUID, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount domain.Account, err error)
	ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID uuid.UUID) ([]domain.Account, error)
}
