package main

import (
	"context"

	"github.com/petrostrak/agile-transfer/internal/data"
	"github.com/shopspring/decimal"
)

type TransferTxParams struct {
	SourceAccountID  int64           `json:"source_account_id"`
	TargetAccountID  int64           `json:"target_account_id"`
	SourceBalance    decimal.Decimal `json:"source_balance"`
	AmountToTransfer decimal.Decimal `json:"amount_to_transfer"`
	SourceCurrency   string          `json:"source_currency"`
	TargetCurrency   string          `json:"target_currency"`
}

type TransferTxResult struct {
	data.Transfer `json:"transfer"`
	SourceAccount data.Account `json:"source_account"`
	TargetAccount data.Account `json:"target_account"`
}

func (app *application) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult

	err := app.models.Transfers.ExecTx(ctx, func() error {
		var err error

		if arg.SourceAccountID == arg.TargetAccountID {
			return ErrIdenticalAccount
		}

		if arg.SourceCurrency != arg.TargetCurrency {
			convertedAmount, err := app.currencyConvertion(arg.SourceCurrency, arg.TargetCurrency, arg.AmountToTransfer)
			if err != nil {
				return ErrCurrencyConvertion
			}

			arg.AmountToTransfer = convertedAmount
		}

		if arg.SourceBalance.LessThan(arg.AmountToTransfer) {
			return ErrInsufficientBalance
		}

		result.SourceAccount, result.TargetAccount, err = app.models.Accounts.AddMoney(arg.SourceAccountID, arg.AmountToTransfer.Neg(), arg.TargetAccountID, arg.AmountToTransfer)
		if err != nil {
			return err
		}

		trasfer := data.Transfer{
			SourceAccountID: arg.SourceAccountID,
			TargetAccountID: arg.TargetAccountID,
			Amount:          arg.AmountToTransfer,
			Currency:        arg.TargetCurrency,
		}
		result.Transfer, err = app.models.Transfers.Insert(trasfer)
		if err != nil {
			return err
		}

		return nil
	})

	return &result, err
}
