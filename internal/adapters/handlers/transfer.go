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
	svc services.TransferService
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
	}

	accounts, err := t.ValidAccounts(ctx, input.SourceAccountID, input.TargetAccountID)
	if err != nil {
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

	result, err := t.svc.TransferTx(ctx, arg)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"transaction": result}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (t *TransferHandler) ValidAccounts(ctx context.Context, sourceAccountID, targetAccountID uuid.UUID) ([]domain.Account, error) {
	return t.svc.ValidateAccounts(ctx, sourceAccountID, targetAccountID)
}

func (t *TransferHandler) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := t.svc.GetAll()
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

	err := t.svc.ExecTx(ctx, func() error {
		var err error

		if arg.SourceAccountID == arg.TargetAccountID {
			return utils.ErrIdenticalAccount
		}

		if arg.SourceCurrency != arg.TargetCurrency {
			convertedAmount, err := utils.CurrencyConvertion(arg.SourceCurrency, arg.TargetCurrency, arg.AmountToTransfer)
			if err != nil {
				return utils.ErrCurrencyConvertion
			}

			arg.AmountToTransfer = convertedAmount
		}

		if arg.SourceBalance.LessThan(arg.AmountToTransfer) {
			return utils.ErrInsufficientBalance
		}

		result.SourceAccount, result.TargetAccount, err = t.svc.AddMoney(
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

		result.Transfer, err = t.svc.Insert(ctx, trasfer)
		if err != nil {
			return err
		}

		return nil
	})

	return &result, err
}
