package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/agile-transfer/internal/adapters/repository"
	"github.com/petrostrak/agile-transfer/internal/core/services"
)

func Test_AccountHandlers(t *testing.T) {
	store := repository.NewPostgressRepository()
	accountService := services.NewAccountService(store.AccountRepository)
	transferService := services.NewTransferService(store.TransferRepository)
	accountHandler := NewAccountHandler(*accountService)
	transferHandler := NewTransferHandler(*transferService)

	testCases := []struct {
		name               string
		method             string
		paramID            string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{"getAllAccounts", "GET", "", accountHandler.GetAllAccounts, http.StatusOK},
		{"getAllAccounts-Invalid", "DELETE", "", accountHandler.GetAllAccounts, http.StatusMethodNotAllowed},
		{"createAccount", "POST", "", accountHandler.CreateAccount, http.StatusCreated},
		{"createAccount-Invalid", "PUT", "", accountHandler.CreateAccount, http.StatusMethodNotAllowed},
		{"getAccount", "GET", "121f03cd-ce8c-447d-8747-fb8cb7aa3a53", accountHandler.GetAccount, http.StatusOK},
		{"getAccount-Invalid", "GET", "121f03cd-ce8c-447d-8747-fb8cb7aa3a52", accountHandler.GetAccount, http.StatusMethodNotAllowed},
		{"updateAccount", "PATCH", "121f03cd-ce8c-447d-8747-fb8cb7aa3a52", accountHandler.UpdateAccount, http.StatusOK},
		{"updateAccount", "PATCH", "121f03cd-ce8c-447d-8747-fb8cb7aq3a52", accountHandler.UpdateAccount, http.StatusMethodNotAllowed},
		{"deleteAccount", "DELETE", "121f03cd-ce8c-447d-8747-fb8cb7aa3a52", accountHandler.DeleteAccount, http.StatusOK},
		{"deleteAccount", "DELETE", "121f039d-ce8c-447d-8747-fb8cb7aa3a52", accountHandler.DeleteAccount, http.StatusMethodNotAllowed},
		{"createTransfer", "POST", "", transferHandler.CreateTransfer, http.StatusCreated},
		{"createTransfer", "GET", "", transferHandler.CreateTransfer, http.StatusMethodNotAllowed},
		{"getAllTransfers", "GET", "", transferHandler.GetAllTransfers, http.StatusOK},
		{"getAllTransfers", "PUT", "", transferHandler.GetAllTransfers, http.StatusMethodNotAllowed},
	}

	for _, tt := range testCases {
		var req *http.Request
		req, _ = http.NewRequest(tt.method, "/", nil)

		if tt.paramID == "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.paramID)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tt.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("%s: wrong status returned; expected %d but got %d", tt.name, tt.expectedStatusCode, rr.Code)
		}
	}
}
