package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
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
)

func TestMain(m *testing.M) {

	code := m.Run()
	os.Exit(code)
}
