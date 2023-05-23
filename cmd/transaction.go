package cmd

import (
	"context"

	"github.com/shopspring/decimal"
)

type accountService struct {
	accountRepo AccountRepository
}

func NewAccountService(accountRepo AccountRepository) AccountService {
	return &accountService{
		accountRepo,
	}
}

func (a *accountService) Insert(acc *Account) error {
	return a.accountRepo.Insert(acc)
}

func (a *accountService) Get(id int64) (*Account, error) {
	return a.accountRepo.Get(id)
}

func (a *accountService) Update(account *Account) error {
	return a.accountRepo.Update(account)
}

func (a *accountService) Delete(id int64) error {
	return a.accountRepo.Delete(id)
}

func (a *accountService) GetAll(ctx context.Context) ([]Account, error) {
	return a.accountRepo.GetAll(ctx)
}

func (a *accountService) AddAccountBalance(ctx context.Context, id int64, amount decimal.Decimal) (Account, error) {
	return a.accountRepo.AddAccountBalance(ctx, id, amount)
}

func (a *accountService) AddMoney(ctx context.Context, sourceAccountID int64, sourceAccountAmount decimal.Decimal, targetAccountID int64, targetAccountAmount decimal.Decimal) (sourceAccount, targetAccount Account, err error) {
	return a.accountRepo.AddMoney(ctx, sourceAccountID, sourceAccountAmount, targetAccountID, targetAccountAmount)
}

func (a *accountService) ValidateAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]Account, error) {
	return a.accountRepo.ValidateAccounts(ctx, sourceAccountID, targetAccountID)
}

type transferService struct {
	transferRepo TransferRepository
}

func NewTransferService(transferRepo TransferRepository) TransferService {
	return &transferService{
		transferRepo,
	}
}

func (t *transferService) Insert(ctx context.Context, tx Transfer) (Transfer, error) {
	return t.transferRepo.Insert(ctx, tx)
}

func (t *transferService) Get(id int64) (*Transfer, error) {
	return t.transferRepo.Get(id)
}

func (t *transferService) GetAll() ([]Transfer, error) {
	return t.transferRepo.GetAll()
}

func (t *transferService) ExecTx(ctx context.Context, fn func() error) error {
	return t.transferRepo.ExecTx(ctx, fn)
}
func (t *transferService) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	return t.transferRepo.TransferTx(ctx, arg)
}
