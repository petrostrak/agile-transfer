package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", app.getAllAccounts)
		r.Post("/", app.createAccount)
		r.Get("/{id}", app.getAccount)
		r.Patch("/{id}", app.updateAccount)
		r.Delete("/{id}", app.deleteAccount)
	})
	r.Post("/transfer", app.createTransfer)
	r.Get("/transactions", app.getAllTransfers)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
