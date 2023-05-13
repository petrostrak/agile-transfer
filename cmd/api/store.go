package main

import (
	"errors"

	"github.com/petrostrak/agile-transfer/internal/data"
)

type TransferTxParams struct {
	SourceAccountID int64   `json:"source_account_id"`
	TargetAccountID int64   `json:"target_account_id"`
	Amount          float64 `json:"amount"`
}

type TransferTxResult struct {
	data.Transfer `json:"transfer"`
	SourceAccount data.Account `json:"source_account"`
	TargetAccount data.Account `json:"target_account"`
}

func (app *application) TransferTx(arg TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult
	var err error
	trasfer := data.Transfer{
		SourceAccountID: arg.SourceAccountID,
		TargetAccountID: arg.TargetAccountID,
		Amount:          arg.Amount,
	}
	result.Transfer, err = app.models.Transfers.Insert(trasfer)
	if err != nil {
		return &result, err
	}

	sourceAccount, err := app.models.Accounts.Get(arg.SourceAccountID)
	if err != nil {
		return nil, err
	}

	if sourceAccount.Balance < arg.Amount {
		return nil, errors.New("insufficient balance")
	} else {
		result.SourceAccount, result.TargetAccount, err = app.models.Accounts.AddMoney(arg.SourceAccountID, -arg.Amount, arg.TargetAccountID, arg.Amount)
	}

	return &result, err
}
