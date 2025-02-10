# Balance Service

This project is a RESTful service developed using **Golang** and **PostgreSQL** for processing and saving incoming requests from third-party providers. The service provides two key routes for managing user account balances.

---

## Features
- **Transaction Processing:** Handle incoming POST requests to update user balances based on transactions.
- **Balance Retrieval:** Query the current balance of a specified user.
- **Concurrency-Safe:** Ensures that transactions are processed safely and efficiently.
- **Idempotency:** Prevents duplicate transaction processing using unique transaction IDs.
- **Transaction Rollbacks:** Handle partial failures and ensure atomicity.
- **Detailed Logging:** Logging is implemented for easy debugging.
- **Testing:** Unit tests are implemented to ensure service works well.
- **Dockerized Deployment:** Easily deployable using Docker and Docker Compose.

---

## Requirements
- **Golang:** 1.18 or higher
- **PostgreSQL:** Database for persisting user data and transactions
- **Docker & Docker Compose:** For easy deployment

---

## Getting Started
### 1. Clone the repository:
```bash
git clone https://github.com/dapoadeleke/balance-service.git
cd balance-service
```

### 2. Build and run using Docker Compose:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:8080`.

### 3. Predefined Users:
Upon startup, the following users will be available:
- User 1
- User 2
- User 3

These users are initialized with zero balance for testing purposes.

---

## API Endpoints

### `POST /user/{userId}/transaction`
Processes a transaction and updates the user's balance.

#### Request Example:
```http
POST /user/1/transaction HTTP/1.1
Host: localhost:8080
Source-Type: game
Content-Type: application/json
```

#### Request Body:
```json
{
  "state": "win",  
  "amount": "10.15",  
  "transactionId": "unique-transaction-id"
}
```

#### Headers:
- `Source-Type`: Indicates the type of transaction (e.g., `game`, `server`, or `payment`).

#### Possible `state` values:
- `win`: Increases the user's balance.
- `lose`: Decreases the user's balance.

#### Response Example:
- **Success:** `200 OK`
- **Error:** Appropriate HTTP status code (e.g., `400 Bad Request` for validation errors)

---

### `GET /user/{userId}/balance`
Retrieves the current balance of a user.

#### Request Example:
```http
GET /user/1/balance HTTP/1.1
Host: localhost:8080
```

#### Response Example:
```json
{
  "userId": 1,
  "balance": "9.25"
}
```

- `balance` is a string rounded to 2 decimal places.

---

### Run Tests:
```bash
go test -v ./...
```

