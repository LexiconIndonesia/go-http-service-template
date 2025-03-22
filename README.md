# Go HTTP Service Template

A modern, production-ready Go HTTP service template following best practices.

## Features

- **Clean Architecture**: Organized with dependency injection pattern
- **Robust Error Handling**: Contextual errors with proper propagation
- **Graceful Shutdown**: Proper signal handling for clean server shutdown
- **Structured Logging**: Using zerolog for performant structured logging
- **Database Integration**: PostgreSQL integration with connection pooling
- **Input Validation**: Request validation using go-playground/validator
- **Middleware Support**: Configurable middleware chain using Chi
- **Testing Support**: Example tests with mocking capabilities
- **Configuration Management**: Environment-based configuration

## Getting Started

### Prerequisites

- Go 1.24.1+
- PostgreSQL 14+
- SQLC (Install: `brew install sqlc`)

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Write SQL queries in `query.sql`
4. Run `sqlc generate`, the generated files will be in `repository` folder
5. Run `make install`

### Development

For local development with hot reloading:

```
./run-dev.sh
```

### Usage

1. Run `make build`
2. Execute `bin/app`

## Project Structure

```
.
├── common/            # Common utilities and models
│   ├── db/            # Database access layer
│   ├── models/        # Domain models
│   └── utils/         # Utility functions
├── middlewares/       # HTTP middleware
├── migrations/        # Database migrations
├── module/            # Feature modules (resources)
├── repository/        # Database repositories (generated)
└── main.go            # Application entry point
```

## Best Practices

This template follows Go best practices including:

- Explicit error handling with context
- Dependency injection instead of global state
- Proper context propagation
- Structured logging
- Graceful shutdowns
- Clean separation of concerns

---

© 2024
