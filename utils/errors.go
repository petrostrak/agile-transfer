package utils

import (
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	ErrIdenticalAccount    = errors.New("source and target account are the same")
	ErrCurrencyConvertion  = errors.New("could not convert currency")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidIDParam      = errors.New("invalid id parameter")
	ErrEmptyBody           = errors.New("body must not be empty")
	ErrBadJSON             = errors.New("body contains badly-formed JSON")
	ErrSingleJSON          = errors.New("body must only contain a single JSON value")

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
)

func LogError(r *http.Request, err error) {
	logger.Println(err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}
