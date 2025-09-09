# Project Structure

## Directory Organization

The project follows Go's standard project layout with clean architecture principles:

```
dfood/
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── api/               # API layer
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── routes/        # Route definitions
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── handler/           # Business logic handlers
│   ├── models/            # Data models and structs
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic layer
│   └── utils/             # Utility functions
├── pkg/                   # Public library code
│   ├── errors/            # Error handling utilities
│   └── logger/            # Logging utilities
├── config/                # Configuration files
│   ├── config.dev.yaml
│   ├── config.staging.yaml
│   └── config.production.yaml
├── api/                   # API documentation (empty)
├── deploy/                # Deployment configurations (empty)
└── scripts/               # Build and deployment scripts (empty)
```

## Architecture Layers

### Repository Layer (`internal/repository/`)
- Defines interfaces for data access
- Implements database operations
- Follows repository pattern

### Service Layer (`internal/service/`)
- Contains business logic
- Orchestrates repository calls
- Handles authentication and authorization

### API Layer (`internal/api/`)
- HTTP handlers and routing
- Request/response processing
- Middleware for cross-cutting concerns

### Models (`internal/models/`)
- Data structures and DTOs
- Database entity definitions
- Request/response models

## Naming Conventions
- Use snake_case for file names
- Use PascalCase for exported types and functions
- Use camelCase for unexported variables and functions
- Interface names should describe behavior (e.g., `UserRepository`)
- Repository implementations end with `Repository` (e.g., `UserRepository`)
- Service implementations end with `Service` (e.g., `AuthService`)

## File Organization Rules
- Keep related functionality in the same package
- Use interfaces to define contracts between layers
- Place shared utilities in `pkg/` for reusable code
- Keep application-specific code in `internal/`
- Configuration files use environment-specific naming