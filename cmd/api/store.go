package main

import (
	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

type TransferTxParams struct {
	SourceAccountID int64           `json:"source_account_id"`
	TargetAccountID int64           `json:"target_account_id"`
	Amount          decimal.Decimal `json:"amount"`
}

type TransferTxResult struct {
	data.Transfer `json:"transfer"`
	SourceAccount data.Account `json:"source_account"`
	TargetAccount data.Account `json:"target_account"`
}

func (app *application) TransferTx(arg TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult

	sourceAccount, err := app.models.Accounts.Get(arg.SourceAccountID)
	if err != nil {
		return nil, err
	}

	targetAccount, err := app.models.Accounts.Get(arg.TargetAccountID)
	if err != nil {
		return nil, err
	}

	if sourceAccount.ID == targetAccount.ID {
		return nil, ErrIdenticalAccount
	}

	if sourceAccount.Currency != targetAccount.Currency {
		convertedAmount, err := app.currencyConvertion(sourceAccount.Currency, targetAccount.Currency, arg.Amount)
		if err != nil {
			return nil, ErrCurrencyConvertion
		}

		arg.Amount = convertedAmount
	}

	if sourceAccount.Balance.LessThan(arg.Amount) {
		return nil, ErrInsufficientBalance
	}

	result.SourceAccount, result.TargetAccount, err = app.models.Accounts.AddMoney(arg.SourceAccountID, arg.Amount.Neg(), arg.TargetAccountID, arg.Amount)
	if err != nil {
		return &result, err
	}

	trasfer := data.Transfer{
		SourceAccountID: arg.SourceAccountID,
		TargetAccountID: arg.TargetAccountID,
		Amount:          arg.Amount,
		Currency:        targetAccount.Currency,
	}
	result.Transfer, err = app.models.Transfers.Insert(trasfer)
	if err != nil {
		return &result, err
	}

	return &result, err
}
