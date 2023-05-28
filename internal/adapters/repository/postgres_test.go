package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/petrostrak/agile-transfer/internal/core/domain"
	"github.com/shopspring/decimal"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbName   = "agile_transfer_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var (
	resource *dockertest.Resource
	pool     *dockertest.Pool
	testDB   *sql.DB
	testRepo PostgresRepository
)

var (
	testID = uuid.UUID{}
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

	testRepo = PostgresRepository{
		&AccountRepository{testDB},
		&TransferRepository{testDB},
	}

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

func Test_PostgresDBRepoInsertAccount(t *testing.T) {
	testAccount := domain.Account{
		Balance:   decimal.NewFromInt(150000),
		Currency:  "EUR",
		CreatedAt: time.Now(),
	}

	err := testRepo.AccountRepository.Insert(&testAccount)
	if err != nil {
		t.Errorf("insert account returned an error: %s", err)
	}
}

func Test_PostgresDBRepoGetAllAccounts(t *testing.T) {
	accounts, err := testRepo.AccountRepository.GetAll(context.Background())
	if err != nil {
		t.Errorf("all accounts report an errorL %s", err)
	}
	testID = accounts[0].ID

	if len(accounts) != 1 {
		t.Errorf("all accounts report wrong size; expected 1, but got %d", len(accounts))
	}

	testAccount := domain.Account{
		Balance:   decimal.NewFromInt(250000),
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	_ = testRepo.AccountRepository.Insert(&testAccount)

	accounts, err = testRepo.AccountRepository.GetAll(context.Background())
	if err != nil {
		t.Errorf("all accounts report an errorL %s", err)
	}

	if len(accounts) != 2 {
		t.Errorf("all accounts report wrong size after insert; expected 2, but got %d", len(accounts))
	}
}

func Test_PostgresDBRepoGetAccount(t *testing.T) {
	account, err := testRepo.AccountRepository.Get(testID)
	if err != nil {
		t.Errorf("error getting account by id: %s", err)
	}

	if account.Currency != "EUR" {
		t.Errorf("wrong account currency returned. expected 'EUR' but got %s", account.Currency)
	}

	if !account.Balance.Equal(decimal.NewFromInt(150000)) {
		t.Errorf("wrong account balance returned. expected 150000 but got %v", account.Balance)
	}
}
