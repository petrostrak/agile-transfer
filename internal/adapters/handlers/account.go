package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/petrostrak/agile-transfer/internal/adapters/repository"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/petrostrak/agile-transfer/internal/core/services"
	"github.com/petrostrak/agile-transfer/utils"
	"github.com/shopspring/decimal"
)

type AccountHandler struct {
	svc services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService,
	}
}

func (a *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Balance  decimal.Decimal `json:"balance"`
		Currency string          `json:"currency"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}

	account := &domain.Account{
		Balance:  input.Balance,
		Currency: input.Currency,
	}

	err = a.svc.Insert(account)
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

func (a *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	account, err := a.svc.Get(id)
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
		ID        uuid.UUID       `json:"id"`
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

func (a *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	account, err := a.svc.Get(id)
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

	err = a.svc.Update(account)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"account": account}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	err := a.svc.Delete(id)
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

func (a *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	accounts, err := a.svc.GetAll(ctx)
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
			ID        uuid.UUID       `json:"id"`
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
