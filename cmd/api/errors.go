package main

import (
	"errors"
	"net/http"
)

var (
	ErrIdenticalAccount    = errors.New("source and target account are the same")
	ErrCurrencyConvertion  = errors.New("could not convert currency")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidIDParam      = errors.New("invalid id parameter")
	ErrEmptyBody           = errors.New("body must not be empty")
	ErrBadJSON             = errors.New("body contains badly-formed JSON")
	ErrSingleJSON          = errors.New("body must only contain a single JSON value")
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
