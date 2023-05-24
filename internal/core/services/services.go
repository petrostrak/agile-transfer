package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/petrostrak/agile-transfer/internal/core/ports"
	"github.com/petrostrak/agile-transfer/utils"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) *AccountService {
	return &AccountService{
		repo,
	}
}

func (a *AccountService) Insert(acc *domain.Account) error {
	return a.repo.Insert(acc)
}

func (a *AccountService) Get(id uuid.UUID) (*domain.Account, error) {
	return a.repo.Get(id)
}

func (a *AccountService) Update(account *domain.Account) error {
	return a.repo.Update(account)
}

func (a *AccountService) Delete(id uuid.UUID) error {
	return a.repo.Delete(id)
}

func (a *AccountService) GetAll(ctx context.Context) ([]domain.Account, error) {
	return a.repo.GetAll(ctx)
}

type TransferService struct {
	repo ports.TransferRepository
}

func NewTransferService(repo ports.TransferRepository) *TransferService {
	return &TransferService{
		repo,
	}
}

func (t *TransferService) Insert(ctx context.Context, tx domain.Transfer) (domain.Transfer, error) {
	return t.repo.Insert(ctx, tx)
}

func (t *TransferService) Get(id uuid.UUID) (*domain.Transfer, error) {
	return t.repo.Get(id)
}

func (t *TransferService) GetAll() ([]domain.Transfer, error) {
	return t.repo.GetAll()
}

func (t *TransferService) ExecTx(ctx context.Context, fn func() error) error {
	return t.repo.ExecTx(ctx, fn)
}
func (t *TransferService) TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error) {
	if arg.SourceAccountID == arg.TargetAccountID {
		return nil, utils.ErrIdenticalAccount
	}
	if arg.SourceCurrency != arg.TargetCurrency {
		convertedAmount, err := utils.CurrencyConvertion(arg.SourceCurrency, arg.TargetCurrency, arg.AmountToTransfer)
		if err != nil {
			return nil, utils.ErrCurrencyConvertion
		}
		arg.AmountToTransfer = convertedAmount
	}
	if arg.SourceBalance.LessThan(arg.AmountToTransfer) {
		return nil, utils.ErrInsufficientBalance
	}
	return t.repo.TransferTx(ctx, arg)
}

func (t *TransferService) AddAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (domain.Account, error) {
	return t.repo.AddAccountBalance(ctx, id, amount)
}

func (t *TransferService) AddMoney(ctx context.Context, sourceAccountID uuid.UUID, sourceAccountAmount decimal.Decimal, targetAccountID uuid.UUID, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount domain.Account, err error) {
	return t.repo.AddMoney(ctx, sourceAccountID, sourceAccountAmount, targetAccountID, targetAccountAmount)
}

func (t *TransferService) ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID uuid.UUID) ([]domain.Account, error) {
	return t.repo.ValidateAccounts(ctx, sourceAccountID, targetAccountID)
}
