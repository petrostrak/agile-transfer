package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

func (app *application) createAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Balance  decimal.Decimal `json:"balance"`
		Currency string          `json:"currency"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	account := &data.Account{
		Balance:  input.Balance,
		Currency: input.Currency,
	}

	err = app.models.Accounts.Insert(account)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/accounts/%d", account.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"account": account}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAccount(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	account, err := app.models.Accounts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var acc struct {
		ID        int64           `json:"id"`
		Balance   decimal.Decimal `json:"balance"`
		Currency  string          `json:"currency"`
		CreatedAt string          `json:"created_at"`
	}
	acc.ID = account.ID
	acc.Balance = account.Balance
	acc.Currency = account.Currency
	acc.CreatedAt = humanDate(account.CreatedAt)

	err = app.writeJSON(w, http.StatusOK, envelope{"account": acc}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	account, err := app.models.Accounts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Balance  *decimal.Decimal `json:"balance"`
		Currency *string          `json:"currency"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	if input.Balance != nil {
		account.Balance = *input.Balance
	}
	if input.Currency != nil {
		account.Currency = *input.Currency
	}

	err = app.models.Accounts.Update(account)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"account": account}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Accounts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "account successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := app.models.Accounts.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var accs []any
	for _, account := range accounts {
		var acc struct {
			ID        int64           `json:"id"`
			Balance   decimal.Decimal `json:"balance"`
			Currency  string          `json:"currency"`
			CreatedAt string          `json:"created_at"`
		}
		acc.ID = account.ID
		acc.Balance = account.Balance
		acc.Currency = account.Currency
		acc.CreatedAt = humanDate(account.CreatedAt)
		accs = append(accs, acc)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"accounts": accs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
