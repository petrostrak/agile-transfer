package main

import (
	"errors"
	"net/http"

	"github.com/disiqueira/gocurrency"
	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

func (app *application) createTransfer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SourceAccountID int64   `json:"source_account_id"`
		TargetAccountID int64   `json:"target_account_id"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
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
		Currency:        input.Currency,
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

func (app *application) currencyConvertion(sourceCurrency, targetCurrency string, amount float64) (decimal.Decimal, error) {
	sourceCur := gocurrency.NewCurrency(sourceCurrency)
	targetCur := gocurrency.NewCurrency(targetCurrency)
	a := decimal.NewFromFloat(amount)

	resultAmount, err := gocurrency.ConvertCurrency(sourceCur, targetCur, a)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return resultAmount, nil
}
