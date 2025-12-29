# Subscription Service

REST API service for managing user online subscriptions.

## Features

- CRUDL operations for subscriptions (Create, Read, Update, Delete, List)
- Calculate total cost of subscriptions for a selected period with filtering
- PostgreSQL database with migrations
- Structured logging with zap
- Swagger API documentation
- Docker Compose support for easy deployment

## Requirements

- Go 1.24+
- PostgreSQL 15+
- Docker and Docker Compose (optional)

## Configuration

Copy `.env.example` to `.env` and configure the following variables:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscription_db
DB_SSLMODE=disable
```

## Running with Docker Compose

The easiest way to run the service:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Subscription service on port 8080

## Running Locally

1. Start PostgreSQL database:
```bash
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=subscription_db \
  -p 5432:5432 \
  postgres:15-alpine
```

2. Run migrations (automatically on startup):
   The service will automatically create the database schema on startup.

3. Build and run the service:
```bash
go mod download
swag init -g cmd/main.go -o docs
go run cmd/main.go
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Subscriptions

- `POST /subscriptions` - Create a new subscription
- `GET /subscriptions` - List all subscriptions (with pagination: `?limit=100&offset=0`)
- `GET /subscriptions/:id` - Get subscription by ID
- `PUT /subscriptions/:id` - Update subscription
- `DELETE /subscriptions/:id` - Delete subscription
- `POST /subscriptions/total-cost` - Calculate total cost for a period

### Swagger Documentation

After starting the service, Swagger UI is available at:
- http://localhost:8080/swagger/index.html

## Example Requests

### Create Subscription

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Calculate Total Cost

```bash
curl -X POST http://localhost:8080/subscriptions/total-cost \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2025",
    "end_date": "12-2025",
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus"
  }'
```

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models
│   ├── repository/          # Database operations
│   └── service/             # Business logic
├── migrations/              # Database migrations
├── docs/                    # Swagger documentation
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Notes

- User existence is not validated (user management is outside the scope)
- Subscription prices are in rubles (integer values only, no kopecks)
- Date format: "MM-YYYY" (e.g., "07-2025")
- The service automatically runs migrations on startup
