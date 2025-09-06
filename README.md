# Order Food Online - Backend

A Go-based backend service for online food ordering, built with Gorilla Mux.

## Quick Start

1. Clone and build:
   ```bash
   git clone https://github.com/ravip18596/oolio-backend-challenge
   cd oolio-backend-challenge
   go mod download
   go build -o bin/server ./cmd/server
   ./bin/server
   ```

2. Access API at `http://localhost:8080`

## API Endpoints

- `GET /product` - List products
- `GET /product/{id}` - Get product details
- `POST /order` - Place order

## Project Structure

- `api/` - OpenAPI specs
- `cmd/server/` - Main application
- `internal/` - Private application code
  - `handler/` - HTTP handlers
  - `model/` - Data models

## Development

Run tests:
```bash
go test ./...
```

Format code:
```bash
gofmt -w .
```
