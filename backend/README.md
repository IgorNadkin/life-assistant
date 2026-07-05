# Life Assistant Backend - Development Setup

This directory contains the backend API for the Life Assistant application.

## Quick Start

### Using Make (Recommended)

```bash
# Show all available commands
make help

# Start PostgreSQL
make db-up

# Run the application in development mode
make dev

# Stop everything
make db-down
```

### Manual Setup

1. **Start PostgreSQL**
   ```bash
   docker-compose up -d
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the server**
   ```bash
   go run cmd/api/main.go
   ```

Server will be available at `http://localhost:8080`

## Project Structure

See [../README.md](../README.md) for detailed project structure.

## API Documentation

See [../docs/API.md](../docs/API.md) for complete API reference.

## Architecture

See [../docs/ARCHITECTURE.md](../docs/ARCHITECTURE.md) for detailed architecture explanation.

## Database

Database migrations are in `migrations/` directory. They are automatically applied on application startup.

### Manual Migration

To manually run migrations:

```bash
# Using psql
psql -h localhost -U postgres -d life < migrations/001_init.sql
psql -h localhost -U postgres -d life < migrations/002_graph.sql
# ... etc
```

## Testing

### Using Postman

1. Import `../postman_collection.json` into Postman
2. Use the predefined requests to test all endpoints

### Manual Testing

```bash
# Health check
curl http://localhost:8080/api/health

# Create scenario
curl -X POST http://localhost:8080/api/scenario \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Смена паспорта",
    "description": "Процесс смены паспорта РФ",
    "start_node_id": 1
  }'
```

## Environment Variables

Create `.env` file based on `.env.example`:

```bash
APP_PORT=8080
APP_ENV=dev
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/life?sslmode=disable
```

## Development Tips

### Code Organization

- **Clean Architecture**: Keep business logic separate from infrastructure
- **No ORM**: Use `sqlx` for database queries
- **Interfaces First**: Define repository interfaces before implementation
- **Dependency Injection**: Use constructor injection for dependencies

### Common Tasks

```bash
# Format code
make fmt

# Run linter
make lint

# Build production binary
make build

# Run all tests
make test

# Clean build artifacts
make clean
```

## Troubleshooting

### Database Connection Error

```
failed to connect postgres: dial tcp localhost:5432: connection refused
```

**Solution:**
```bash
make db-up
# Wait for container to start
sleep 2
make dev
```

### Port Already in Use

```
listen tcp :8080: bind: address already in use
```

**Solution:**
```bash
# Kill process on port 8080
lsof -ti :8080 | xargs kill -9
```

### Module Not Found

```
no required module provides package ...
```

**Solution:**
```bash
make deps
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
