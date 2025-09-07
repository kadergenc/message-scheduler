- [Message Scheduler](#message-scheduler)
    * [Features](#features)
    * [Technologies Used](#technologies-used)
        + [Backend](#backend)
        + [Database](#database)
        + [Testing & Development](#testing---development)
        + [Infrastructure & DevOps](#infrastructure---devops)
    * [Prerequisites](#prerequisites)
    * [Installation](#installation)
        + [1. Clone the Repository](#1-clone-the-repository)
        + [2. Install Dependencies](#2-install-dependencies)
        + [3. Set Up Database](#3-set-up-database)
    * [Configuration](#configuration)
        + [1. Database Configuration](#1-database-configuration)
        + [2. Application Configuration](#2-application-configuration)
    * [Usage](#usage)
        + [Starting the Service](#starting-the-service)
        + [API Endpoints](#api-endpoints)
            - [Start Message Processing](#start-message-processing)
            - [Stop Message Processing](#stop-message-processing)
            - [Health Check](#health-check)
            - [API Documentation](#api-documentation)
    * [Webhook Integration](#webhook-integration)
    * [Testing](#testing)
        + [Run All Tests](#run-all-tests)
        + [Run Tests with Coverage](#run-tests-with-coverage)
        + [Run Application Layer Tests](#run-application-layer-tests)
        + [Generate Mocks](#generate-mocks)
    * [Project Structure](#project-structure)
    * [Architecture](#architecture)
        + [Key Components](#key-components)
    * [Development](#development)
        + [Adding New Features](#adding-new-features)
        + [Code Standards](#code-standards)
    * [Monitoring](#monitoring)
  
# Message Scheduler

A robust Go-based message scheduling service that processes and sends messages via webhooks with automatic retry mechanisms and comprehensive logging.

## Features

- **Scheduled Message Processing**: Automatically processes unsent messages at configurable intervals
- **Webhook Integration**: Sends messages through HTTP webhooks with timeout configuration
- **Database Persistence**: PostgreSQL storage with GORM ORM for reliable data management
- **REST API**: RESTful endpoints for message management and scheduler control
- **Comprehensive Logging**: Structured logging with zerolog for monitoring and debugging
- **Error Handling**: Graceful error handling with automatic status updates
- **Unit Testing**: 96.6% test coverage with testify and mockery
- **Graceful Shutdown**: Clean shutdown handling for production environments

## Technologies Used

### Backend
- **[Go 1.23+](https://golang.org/)** - Primary programming language
- **[Fiber v2](https://gofiber.io/)** - Fast HTTP web framework built on Fasthttp
- **[GORM](https://gorm.io/)** - Go ORM library for database operations
- **[Zerolog](https://github.com/rs/zerolog)** - Structured, high-performance logging

### Database
- **[PostgreSQL 12+](https://www.postgresql.org/)** - Primary relational database

### Testing & Development
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit with assertions and mocks
- **[Mockery](https://github.com/vektra/mockery)** - Mock code generator for Go interfaces
- **[Swagger](https://swagger.io/)** - API documentation and testing interface

### Infrastructure & DevOps
- **[Docker](https://www.docker.com/)** - Containerization platform
- **[Docker Compose](https://docs.docker.com/compose/)** - Multi-container application orchestration

## Prerequisites

Before running this application, ensure you have:

- **Go 1.23+** installed
- **PostgreSQL 12+** running
- **Git** for version control

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd message-scheduler
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Database

Create a PostgreSQL database and run the schema

First, make sure your Docker daemon is running then run the following command to create the database and apply the schema:

```bash
cd local
docker compose up --build
```

This will spin up a PostgreSQL Docker container and automatically execute init.sql inside it to populate your database with data.

## Configuration

### 1. Database Configuration

Create `config/pg-credentials.json`:

```json
{
  "username": "your_db_user",
  "password": "your_db_password"
}
```

### 2. Application Configuration

Update `config/config.json`:

```json
{
  "port": "8081",
  "appName": "message-scheduler",
  "webhook": {
    "host": "http://your-webhook-url:8000/webhook",
    "timeout": 3000
  },
  "postgres": {
    "writeHost": "localhost",
    "writePort": "5432",
    "readHost": "localhost",
    "readPort": "5432",
    "dbname": "messaging"
  }
}
```

## Usage

### Starting the Service

```bash
# Run with default config
go run main.go

# Run with custom config paths
go run main.go -config ./config/config.json -pg ./config/pg-credentials.json
```

The service will start on the configured port (default: 8081) and display:
```
Message Scheduler started successfully. Call /start-send-message to begin processing.
```

### API Endpoints

#### Start Message Processing
```http
POST /start-send-message
```
Starts the automatic message processing scheduler.

#### Stop Message Processing
```http
POST /stop-send-message
```
Stops the message processing scheduler.

#### Health Check
```http
GET /health
```
Returns service health status.

#### API Documentation
```http
GET /swagger/*
```
Access Swagger UI for complete API documentation.

## Webhook Integration

The service sends messages to your configured webhook endpoint with the following payload:

```json
{
  "to": "+1234567890",
  "content": "Your message content here"
}
```

Expected webhook response:
```json
{
  "message": "Message sent successfully",
  "messageId": "webhook-generated-id"
}
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Run Application Layer Tests
```bash
go test ./internal/application -v -cover
```

### Generate Mocks
```bash
mockery --config .mockery.yaml
```

## Project Structure

```
message-scheduler/
├── cmd/                    # Application entrypoints
├── config/                 # Configuration files
│   ├── config.json        # Main application config
│   ├── config-dev.json    # Development config
│   └── pg-credentials.json # Database credentials
├── docs/                   # API documentation
├── internal/
│   ├── application/        # Business logic layer
│   │   ├── message-send-service.go
│   │   └── message_send_service_test.go
│   ├── domain/            # Domain entities and types
│   │   ├── entity/        # Domain entities
│   │   └── types/         # Domain types
│   ├── infra/             # Infrastructure layer
│   │   ├── client/        # External service clients
│   │   ├── database/      # Database configuration
│   │   ├── repository/    # Data access layer
│   │   ├── scheduler/     # Job scheduler
│   │   └── server/        # HTTP server
│   └── port/              # Interfaces/contracts
├── local/                 # Local development setup
│   └── docker-compose.yaml
├── log/                   # Logging configuration
├── mocks/                 # Generated mocks for testing
├── main.go               # Application entry point
├── go.mod               # Go module definition
└── README.md           # This file
```

## Architecture

The application follows Clean Architecture principles:

- **Domain Layer**: Core business entities and types
- **Application Layer**: Business logic and use cases
- **Infrastructure Layer**: External concerns (database, HTTP, webhooks)
- **Port Layer**: Interfaces defining contracts between layers

### Key Components

- **MessageSendService**: Core service for message processing
- **PostgresRepository**: Database operations with GORM
- **WebhookClient**: HTTP client for webhook communication  
- **SimpleScheduler**: Job scheduling with configurable intervals
- **Fiber Server**: HTTP REST API server

## Development

### Adding New Features

1. Define interfaces in `internal/port/`
2. Implement business logic in `internal/application/`
3. Add infrastructure in `internal/infra/`
4. Write comprehensive tests
5. Update API documentation

### Code Standards

- Follow Go conventions and best practices
- Maintain test coverage above 90%
- Use structured logging with appropriate log levels
- Handle errors gracefully with proper context
- Document public APIs with comments

## Monitoring

The application provides extensive logging for monitoring:

- **Info Level**: Normal operations and state changes
- **Error Level**: Failed operations with context
- **Debug Level**: Detailed execution information

Log format is JSON for easy parsing and analysis.

---