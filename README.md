# Auto Message Dispatcher

An automatic message sending system that processes and sends messages from a database at scheduled intervals.

## Requirements

- Go 1.20+
- MongoDB
- Redis
- Docker & Docker Compose (optional)

## Installation

### With Docker

```bash
# Start all services
docker-compose up -d

# API available at: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger
```

### Manual Setup

```bash
# Download dependencies
go mod download

# Create .env file
cp .env.example .env
# Edit .env and add your configuration values

# Start MongoDB (local)
mongod

# Start Redis (optional)
redis-server

# Run the application
go run cmd/server/main.go
```

## API Endpoints

### Scheduler Control
- `POST /api/scheduler/start` - Start automatic message sending
- `POST /api/scheduler/stop` - Stop automatic message sending
- `GET /api/scheduler/status` - Check scheduler status

### Message Operations
- `GET /api/messages/sent` - List sent messages
- `POST /api/messages` - Create new message

## Configuration

Create a `.env` file from `.env.example` and configure your values.

Important variables:
- `WEBHOOK_URL`: Message sending endpoint
- `WEBHOOK_AUTH_KEY`: API authentication key
- `SCHEDULER_INTERVAL`: Message sending interval (default: 2m)
- `SCHEDULER_BATCH_SIZE`: Number of messages per batch (default: 2)
- `SCHEDULER_AUTO_START`: Auto-start on deployment (default: true)

## Features

- Custom scheduler implementation (no external cron packages)
- Automatic message sending on deployment
- Prevents duplicate message sending
- Redis caching for sent messages (bonus feature)
- Retry mechanism with exponential backoff
- Configuration validation on startup
- Swagger/OpenAPI documentation
- Docker support with health checks

## Swagger Documentation

Access API documentation at: `http://localhost:8080/swagger`

## Architecture

The project follows clean architecture principles with clear separation of concerns:

- **Domain**: Core business entities and logic
- **Repository**: Data access layer
- **Service**: Business logic and external integrations
- **Handler**: HTTP request handlers
- **Scheduler**: Custom cron-like implementation
- **Middleware**: Cross-cutting concerns (CORS, logging)

