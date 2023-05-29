// go:build integration
package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/petrostrak/agile-transfer/internal/adapters/repository"
	"github.com/petrostrak/agile-transfer/internal/core/services"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbName   = "agile_transfer_test"
	port     = "5436"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var (
	resource        *dockertest.Resource
	pool            *dockertest.Pool
	testDB          *sql.DB
	testRepo        repository.PostgresRepository
	accountService  *services.AccountService
	transferService *services.TransferService
	accountHandler  *AccountHandler
	transferHandler *TransferHandler
)

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&options)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("postgres", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to DB: %s", err)
	}

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = repository.PostgresRepository{
		&repository.AccountRepository{testDB},
		&repository.TransferRepository{testDB},
	}

	accountService = services.NewAccountService(testRepo.AccountRepository)
	transferService = services.NewTransferService(testRepo.TransferRepository)
	accountHandler = NewAccountHandler(*accountService)
	transferHandler = NewTransferHandler(*transferService)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/init_schema.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_PingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("cannot ping DB")
	}
}

func Test_AccountHandlers(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		json               string
		paramID            string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{"getAllAccounts", "GET", "", "", accountHandler.GetAllAccounts, http.StatusOK},
		{
			"createAccount",
			"POST",
			`{"balance": 150000,"currency": "EUR"}`,
			"",
			accountHandler.CreateAccount,
			http.StatusCreated,
		},
		{"getAccount", "GET", "", "604f02b2-4e45-48d6-a952-03a0136e8140", accountHandler.GetAccount, http.StatusOK},
		{"getAccount-Invalid", "", "GET", "121f03cd-ce8c-447d-8747-fb8cb7aa3a52", accountHandler.GetAccount, http.StatusMethodNotAllowed},
		{
			"updateAccount",
			"PATCH",
			`{"balance": 1500000}`,
			"604f02b2-4e45-48d6-a952-03a0136e8140",
			accountHandler.UpdateAccount,
			http.StatusOK,
		},
		{"deleteAccount", "DELETE", "", "71376d61-8b6c-4289-b5c4-79cb36add23f", accountHandler.DeleteAccount, http.StatusOK},
		{
			"createTransfer",
			"POST",
			`{"source_account_id": "8fa6c93b-f300-4ef8-9bac-4258caea36db","target_account_id": "604f02b2-4e45-48d6-a952-03a0136e8140","amount": 1000000}`,
			"",
			transferHandler.CreateTransfer,
			http.StatusCreated,
		},
		{"getAllTransfers", "GET", "", "", transferHandler.GetAllTransfers, http.StatusOK},
	}

	for _, tt := range testCases {
		var req *http.Request
		if tt.json == "" {
			req, _ = http.NewRequest(tt.method, "/", nil)
		} else {
			req, _ = http.NewRequest(tt.method, "/", strings.NewReader(tt.json))
		}

		if tt.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tt.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("%s: wrong status returned; expected %d but got %d", tt.name, tt.expectedStatusCode, rr.Code)
		}
	}
}
