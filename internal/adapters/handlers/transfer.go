package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/petrostrak/agile-transfer/internal/adapters/repository"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/petrostrak/agile-transfer/internal/core/services"
	"github.com/petrostrak/agile-transfer/utils"
	"github.com/shopspring/decimal"
)

type TransferHandler struct {
	service services.TransferService
}

func NewTransferHandler(transferService services.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService,
	}
}

func (t *TransferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var input struct {
		SourceAccountID uuid.UUID       `json:"source_account_id"`
		TargetAccountID uuid.UUID       `json:"target_account_id"`
		Amount          decimal.Decimal `json:"amount"`
		Currency        string          `json:"currency"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	accounts, err := t.service.ValidateAccounts(ctx, input.SourceAccountID, input.TargetAccountID)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	arg := domain.TransferTxParams{
		SourceAccountID:  input.SourceAccountID,
		TargetAccountID:  input.TargetAccountID,
		SourceBalance:    accounts[0].Balance,
		SourceCurrency:   accounts[0].Currency,
		TargetCurrency:   accounts[1].Currency,
		AmountToTransfer: input.Amount,
	}

	result, err := t.service.TransferTx(ctx, arg)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"transaction": result}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (t *TransferHandler) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := t.service.GetAll()
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

func (t *TransferHandler) TransferTx(ctx context.Context, arg domain.TransferTxParams) (*domain.TransferTxResult, error) {
	var result domain.TransferTxResult

	err := t.service.ExecTx(ctx, func() error {
		var err error

		result.SourceAccount, result.TargetAccount, err = t.service.AddMoney(
			ctx,
			arg.SourceAccountID,
			arg.AmountToTransfer.Neg(),
			arg.TargetAccountID,
			arg.AmountToTransfer,
		)
		if err != nil {
			return err
		}

		trasfer := domain.Transfer{
			SourceAccountID: arg.SourceAccountID,
			TargetAccountID: arg.TargetAccountID,
			Amount:          arg.AmountToTransfer,
			Currency:        arg.TargetCurrency,
		}

		result.Transfer, err = t.service.Insert(ctx, trasfer)
		if err != nil {
			return err
		}

		return nil
	})

	return &result, err
}
