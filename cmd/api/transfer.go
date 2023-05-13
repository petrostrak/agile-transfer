package main

import (
	"errors"
	"fmt"
	"net/http"

	currconv "github.com/kitloong/go-currency-converter-api/v2"
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

func (app *application) currencyConvertion(from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	// fromCurrency := gocurrency.NewCurrency(from)
	// toCurrency := gocurrency.NewCurrency(to)

	// resultAmount, err := gocurrency.ConvertCurrency(fromCurrency, toCurrency, amount)
	// if err != nil {
	// 	return decimal.Decimal{}, err
	// }

	api := currconv.NewAPI(currconv.Config{
		BaseURL: "https://free.currconv.com",
		Version: "v7",
		APIKey:  "[KEY]",
	})

	convert, err := api.Convert(currconv.ConvertRequest{
		Q: []string{fmt.Sprintf("%s_%s", from, to)},
	})
	if err != nil {
		return decimal.Decimal{}, err
	}

	return decimal.NewFromFloat32(convert.Results[fmt.Sprintf("%s_%s", from, to)].Val), nil
}
