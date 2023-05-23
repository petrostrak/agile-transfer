package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/petrostrak/agile-transfer/cmd"
	"github.com/petrostrak/agile-transfer/repository"
	"github.com/petrostrak/agile-transfer/utils"
	"github.com/shopspring/decimal"
)

type transferHandler struct {
	transferService cmd.TransferService
}

func (t *transferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var input struct {
		SourceAccountID int64           `json:"source_account_id"`
		TargetAccountID int64           `json:"target_account_id"`
		Amount          decimal.Decimal `json:"amount"`
		Currency        string          `json:"currency"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}

	accounts, err := t.ValidAccounts(ctx, input.SourceAccountID, input.TargetAccountID)
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

	result, err := t.transferService.TransferTx(ctx, arg)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"transaction": result}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (t *transferHandler) ValidAccounts(ctx context.Context, sourceAccountID, targetAccountID int64) ([]repository.Account, error) {
	return ValidateAccounts(ctx, sourceAccountID, targetAccountID)
}

func (t *transferHandler) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := t.transferService.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"transfers": transfers}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}
