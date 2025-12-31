# Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

## Setup

Install golang-migrate:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Running Migrations

### Up (apply migrations)
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subscription_db?sslmode=disable" up
```

### Down (rollback migrations)
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subscription_db?sslmode=disable" down
```

### Using environment variables
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/subscription_db?sslmode=disable"
migrate -path migrations -database $DATABASE_URL up
```

## Migration Files

- `001_create_subscriptions_table.up.sql` - Creates subscriptions table
- `001_create_subscriptions_table.down.sql` - Drops subscriptions table

## Docker Compose

When using docker-compose, migrations can be run automatically or manually:

```bash
docker-compose exec app migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/subscription_db?sslmode=disable" up
```

