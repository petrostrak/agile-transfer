package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

func (app *application) createTransfer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

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

	accounts, err := app.validAccounts(ctx, input.SourceAccountID, input.TargetAccountID)
	if err != nil {
		return
	}

	arg := TransferTxParams{
		SourceAccountID:  input.SourceAccountID,
		TargetAccountID:  input.TargetAccountID,
		SourceBalance:    accounts[0].Balance,
		SourceCurrency:   accounts[0].Currency,
		TargetCurrency:   accounts[1].Currency,
		AmountToTransfer: input.Amount,
	}

	result, err := app.TransferTx(ctx, arg)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"transaction": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) validAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]data.Account, error) {
	return app.models.Accounts.ValidateAccounts(ctx, sourceAccountID, targetAccountID)
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
