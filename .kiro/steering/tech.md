# Technology Stack

## Core Technologies
- **Language**: Go 1.24.2
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: SQLite3 with go-sqlite3 driver
- **Authentication**: JWT tokens (golang-jwt/jwt/v5)
- **Configuration**: YAML-based config files
- **Logging**: Custom structured logger

## Key Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/mattn/go-sqlite3` - SQLite database driver
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `gopkg.in/yaml.v3` - YAML configuration parsing
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/go-chi/httprate` - Rate limiting

## Build System
- Standard Go modules (`go.mod`)
- No additional build tools required

## Common Commands
```bash
# Run the application
go run cmd/main.go

# Build the application
go build -o dfood cmd/main.go

# Run tests
go test ./...

# Install dependencies
go mod tidy

# Update dependencies
go mod download
```

## Environment Configuration
- Set `APP_ENV` environment variable to control config file selection
- Supported environments: `dev`, `staging`, `production`
- Config files located in `config/config.{env}.yaml`