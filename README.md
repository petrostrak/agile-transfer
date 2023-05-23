## Agile Transfer
A simple microservice which handles financial transactions written in GO.

To launch the application, mount the root directory of the project and run docker:
```
docker compose up -d
```
and
```
go run ./cmd/api
```

To migrate up
```
migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/agile_transfer?sslmode=disable" -verbose up
```

To migrate down
```
migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/agile_transfer?sslmode=disable" -verbose down
```

While the application is running, we can make requests to add, update, remove accounts and make transactions between them.

*   Create Account (POST) to `localhost:8080/accounts` with request body:

    ```
    {
        "balance": 53000.50,
        "currency": "EUR"
    }
*   Get Account (GET) to `localhost:8080/accounts/{id}`
*   Update Account (PATCH) to `localhost:8080/accounts/{id}` with request body:

    ```
    {
        "balance": 15000,
    }
    ```
*   Delete Account (DELETE) to `localhost:8080/accounts/{id}`
*   Get All Accounts (GET) to `localhost:8080/accounts`
*   Make Transaction (POST) to `localhost:8080/transfer` with request body:
    ```
    {
        "source_account_id": 3,
        "target_account_id": 4,
        "amount": 10
    }
    ```
*   Get All Transactions (GET) to `localhost:8080/transactions`