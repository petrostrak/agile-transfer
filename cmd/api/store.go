package main

import (
	"errors"

	"github.com/petrostrak/agile-transfer/internal/data"
)

type TransferTxParams struct {
	SourceAccountID int64   `json:"source_account_id"`
	TargetAccountID int64   `json:"target_account_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

type TransferTxResult struct {
	data.Transfer `json:"transfer"`
	SourceAccount data.Account `json:"source_account"`
	TargetAccount data.Account `json:"target_account"`
}

func (app *application) TransferTx(arg TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult
	var err error

	sourceAccount, err := app.models.Accounts.Get(arg.SourceAccountID)
	if err != nil {
		return nil, err
	}

	if sourceAccount.Currency != arg.Currency {
		convertedAmount, err := app.currencyConvertion(sourceAccount.Currency, arg.Currency, arg.Amount)
		if err != nil {
			return nil, errors.New("could not convert currency")
		}
		a, err := app.toFloat64(convertedAmount)
		if err != nil {
			return nil, err
		}
		arg.Amount = a
	}

	if sourceAccount.Balance < arg.Amount {
		return nil, errors.New("insufficient balance")
	} else {
		result.SourceAccount, result.TargetAccount, err = app.models.Accounts.AddMoney(arg.SourceAccountID, -arg.Amount, arg.TargetAccountID, arg.Amount)
		if err != nil {
			return &result, err
		}

		trasfer := data.Transfer{
			SourceAccountID: arg.SourceAccountID,
			TargetAccountID: arg.TargetAccountID,
			Amount:          arg.Amount,
			Currency:        sourceAccount.Currency,
		}
		result.Transfer, err = app.models.Transfers.Insert(trasfer)
		if err != nil {
			return &result, err
		}
	}

	return &result, err
}
