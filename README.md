## Agile Transfer
A simple microservice which handles financial transactions written in GO.

To launch the application, mount the root directory of the project and run docker:
```
docker compose up -d
```
and
```
make start
```

To migrate up
```
make migrate-up
```

To migrate down
```
make migrate-down
```

To run tests
```
make coverage
```

To run tests with integration
```
make coverage-integration
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
        "source_account_id": "ac629895-57b4-46f2-bf11-1011fbb015c3",
        "target_account_id": "5531dc5a-4dc2-4e34-97fc-78e4d88d0e22",
        "amount": 15000
    }
    ```
*   Get All Transactions (GET) to `localhost:8080/transactions`