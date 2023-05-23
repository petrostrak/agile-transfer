package services

import (
	"context"

	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/petrostrak/agile-transfer/internal/core/ports"
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

func (a *AccountService) Get(id int64) (*domain.Account, error) {
	return a.repo.Get(id)
}

func (a *AccountService) Update(account *domain.Account) error {
	return a.repo.Update(account)
}

func (a *AccountService) Delete(id int64) error {
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

func (t *TransferService) Get(id int64) (*domain.Transfer, error) {
	return t.repo.Get(id)
}

func (t *TransferService) GetAll() ([]domain.Transfer, error) {
	return t.repo.GetAll()
}

func (t *TransferService) ExecTx(ctx context.Context, fn func() error) error {
	return t.repo.ExecTx(ctx, fn)
}
func (t *TransferService) TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error) {
	return t.repo.TransferTx(ctx, arg)
}

func (t *TransferService) AddAccountBalance(ctx context.Context, id int64, amount decimal.Decimal) (domain.Account, error) {
	return t.repo.AddAccountBalance(ctx, id, amount)
}

func (t *TransferService) AddMoney(ctx context.Context, sourceAccountID int64, sourceAccountAmount decimal.Decimal, targetAccountID int64, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount domain.Account, err error) {
	return t.repo.AddMoney(ctx, sourceAccountID, sourceAccountAmount, targetAccountID, targetAccountAmount)
}

func (t *TransferService) ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]domain.Account, error) {
	return t.repo.ValidateAccounts(ctx, sourceAccountID, targetAccountID)
}
