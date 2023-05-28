start:
	go run main.go
test:
	go test -v ./...
test-integration:
	go test -v -tags=integration ./...
coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
coverage-integration:
	go test -tags=integration ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
migrate-up:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/agile_transfer?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/agile_transfer?sslmode=disable" -verbose down
PHONY: start, coverage, coverage-integration, test, test-integration
