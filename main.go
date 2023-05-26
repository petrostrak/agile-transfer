package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/petrostrak/agile-transfer/internal/adapters/handlers"
	"github.com/petrostrak/agile-transfer/internal/adapters/repository"
	"github.com/petrostrak/agile-transfer/internal/core/services"
)

var (
	accountService  *services.AccountService
	transferService *services.TransferService
	accountHandler  *handlers.AccountHandler
	transferHandler *handlers.TransferHandler
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	store := repository.NewPostgressRepository()
	accountService = services.NewAccountService(store.AccountRepository)
	transferService = services.NewTransferService(store.TransferRepository)
	accountHandler = handlers.NewAccountHandler(*accountService)
	transferHandler = handlers.NewTransferHandler(*transferService)

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", 8080),
		Handler:     Routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
	}

	logger.Printf("starting development server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", accountHandler.GetAllAccounts)
		r.Post("/", accountHandler.CreateAccount)
		r.Get("/{id}", accountHandler.GetAccount)
		r.Patch("/{id}", accountHandler.UpdateAccount)
		r.Delete("/{id}", accountHandler.DeleteAccount)
	})
	r.Post("/transfer", transferHandler.CreateTransfer)
	r.Get("/transactions", transferHandler.GetAllTransfers)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
