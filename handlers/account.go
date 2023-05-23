package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/petrostrak/agile-transfer/cmd"
	"github.com/petrostrak/agile-transfer/repository"
	"github.com/petrostrak/agile-transfer/utils"
	"github.com/shopspring/decimal"
)

type accountHandler struct {
	accountService cmd.AccountService
}

func (a *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Balance  decimal.Decimal `json:"balance"`
		Currency string          `json:"currency"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}

	account := &cmd.Account{
		Balance:  input.Balance,
		Currency: input.Currency,
	}

	err = a.accountService.Insert(account)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/accounts/%d", account.ID))

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"account": account}, headers)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *accountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.NotFoundResponse(w, r)
		return
	}

	account, err := a.accountService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
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
	acc.CreatedAt = utils.HumanDate(account.CreatedAt)

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"account": acc}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *accountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.NotFoundResponse(w, r)
		return
	}

	account, err := a.accountService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Balance  *decimal.Decimal `json:"balance"`
		Currency *string          `json:"currency"`
	}

	err = utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}
	if input.Balance != nil {
		account.Balance = *input.Balance
	}
	if input.Currency != nil {
		account.Currency = *input.Currency
	}

	err = a.accountService.Update(account)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"account": account}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *accountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.NotFoundResponse(w, r)
		return
	}

	err = a.accountService.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "account successfully deleted"}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *accountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	accounts, err := a.accountService.GetAll(ctx)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
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
		acc.CreatedAt = utils.HumanDate(account.CreatedAt)
		accs = append(accs, acc)
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"accounts": accs}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}
