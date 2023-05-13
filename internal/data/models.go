package data

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not Found")

type Models struct {
	Accounts  AccountModel
	Transfers TransferModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Accounts:  AccountModel{DB: db},
		Transfers: TransferModel{DB: db},
	}
}
