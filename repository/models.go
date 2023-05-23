package repository

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not Found")

type Models struct {
	Accounts  accountRepository
	Transfers transferRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Accounts:  accountRepository{DB: db},
		Transfers: transferRepository{DB: db},
	}
}
