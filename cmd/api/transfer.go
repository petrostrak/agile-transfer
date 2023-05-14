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

	err = app.validAccount(w, r, input.SourceAccountID, input.Currency)
	if err != nil {
		return
	}

	err = app.validAccount(w, r, input.TargetAccountID, input.Currency)
	if err != nil {
		return
	}

	arg := TransferTxParams{
		SourceAccountID: input.SourceAccountID,
		TargetAccountID: input.TargetAccountID,
		Amount:          input.Amount,
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

func (app *application) validAccount(w http.ResponseWriter, r *http.Request, accountID int64, currency string) error {
	_, err := app.models.Accounts.Get(accountID)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "One or more of the accounts does not exist")
		return err
	}

	return nil
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
