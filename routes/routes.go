package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/agile-transfer/handlers"
)

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", handlers.GetAllAccounts)
		r.Post("/", handlers.CreateAccount)
		r.Get("/{id}", handlers.GetAccount)
		r.Patch("/{id}", handlers.UpdateAccount)
		r.Delete("/{id}", handlers.DeleteAccount)
	})
	r.Post("/transfer", handlers.CreateTransfer)
	r.Get("/transactions", handlers.GetAllTransfers)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
