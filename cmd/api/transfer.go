package main

import (
	"errors"
	"net/http"

	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

func (app *application) createTransfer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SourceAccountID int64           `json:"source_account_id"`
		TargetAccountID int64           `json:"target_account_id"`
		Amount          decimal.Decimal `json:"amount"`
		Currency        string          `json:"currency"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	sourceAcc, err := app.validAccount(w, r, input.SourceAccountID)
	if err != nil {
		return
	}

	targetAcc, err := app.validAccount(w, r, input.TargetAccountID)
	if err != nil {
		return
	}

	arg := TransferTxParams{
		SourceAccountID:  input.SourceAccountID,
		TargetAccountID:  input.TargetAccountID,
		SourceBalance:    sourceAcc.Balance,
		SourceCurrency:   sourceAcc.Currency,
		TargetCurrency:   targetAcc.Currency,
		AmountToTransfer: input.Amount,
	}

	result, err := app.TransferTx(arg)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"transaction": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) validAccount(w http.ResponseWriter, r *http.Request, accountID int64) (*data.Account, error) {
	account, err := app.models.Accounts.Get(accountID)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "One or more of the accounts does not exist")
		return nil, err
	}

	return account, nil
}

func (app *application) getAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := app.models.Transfers.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transfers": transfers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
