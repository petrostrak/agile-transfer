package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", app.getAllAccounts)
		r.Post("/", app.createAccount)
		r.Get("/{id}", app.getAccount)
		r.Patch("/{id}", app.updateAccount)
		r.Delete("/{id}", app.deleteAccount)
	})
	r.Post("/transfer", app.createTransfer)
	r.Get("/transactions", app.getAllTransfers)

	return r
}
