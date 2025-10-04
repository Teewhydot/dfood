# API Request Flow Documentation

This document explains how API calls are handled in the dfood application, from the initial HTTP request to the final response. We'll use the **User Login** endpoint as an example to demonstrate the complete flow.

## Architecture Overview

The dfood application follows a **Clean Architecture** pattern with clear separation of concerns:

```
HTTP Request → Middleware → Routes → Handlers → Services → Repository → Database
                ↓
HTTP Response ← Error Handler ← Response Builder ← Business Logic ← Data Access ← GORM
```

## Complete Request Flow: POST /api/v1/auth/login

Let's trace a login request through the entire system:

### 1. HTTP Request Arrives
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### 2. Gin Router Processing (`cmd/main.go`)

The request first hits the Gin router configured in `main.go`:

```go
// main.go - Application Entry Point
func main() {
    // ... configuration setup ...
    
    // Dependency injection
    userRepo := repository.NewUserRepository()
    authService := service.NewAuthService(userRepo)
    
    deps := &routes.Dependencies{
        AuthService: authService,
    }
    
    // Setup routes with middleware
    router := routes.SetupRoutes(deps)
    router.Run(":" + fmt.Sprint(cfg.Port))
}
```

### 3. Middleware Chain Processing (`internal/api/routes/routes.go`)

Before reaching the handler, the request passes through multiple middleware layers:

```go
// routes.go - Middleware Chain
func SetupRoutes(deps *Dependencies) *gin.Engine {
    router := gin.New()

    // 1. Request Logging Middleware
    router.Use(middleware.RequestLogger())
    
    // 2. CORS Middleware  
    router.Use(middleware.CORS())
    
    // 3. Recovery Middleware (panic recovery)
    router.Use(gin.Recovery())
    
    // 4. Rate Limiting Middleware
    router.Use(middleware.RateLimitMiddleware(10, time.Minute))
    
    // Route registration
    authHandler := handlers.NewAuthHandler(deps.AuthService)
    v1 := router.Group("/api/v1")
    {
        auth := v1.Group("/auth")
        {
            auth.POST("/login", authHandler.Login) // ← Our endpoint
        }
    }
}
```

#### 3.1 Request Logger Middleware (`internal/api/middleware/logging.go`)
```go
func RequestLogger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        logger.Info("Request processed",
            "client_id", param.ClientIP,
            "method", param.Method,
            "path", param.Path,
            "status_code", param.StatusCode,
            "latency", param.Latency,
        )
        return ""
    })
}
```
**Purpose**: Logs all incoming requests with details like IP, method, path, status code, and latency.

#### 3.2 CORS Middleware (`internal/api/middleware/cors.go`)
```go
func CORS() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
        // ... other CORS headers
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })
}
```
**Purpose**: Handles Cross-Origin Resource Sharing, allowing frontend applications to make requests.

#### 3.3 Rate Limiting Middleware (`internal/api/middleware/rate_limiter.go`)
```go
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
    limiter := httprate.Limit(
        limit,    // 10 requests
        window,   // per minute
        httprate.WithKeyFuncs(httprate.KeyByIP), // per IP
    )
    
    return func(c *gin.Context) {
        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            c.Next()
        })
        limitedHandler := limiter(handler)
        limitedHandler.ServeHTTP(c.Writer, c.Request)
    }
}
```
**Purpose**: Prevents abuse by limiting requests to 10 per minute per IP address.

### 4. Handler Layer (`internal/api/handlers/auth.go`)

After middleware processing, the request reaches the authentication handler:

```go
func (h *AuthHandler) Login(c *gin.Context) {
    // 1. Parse and validate JSON request body
    var loginUser models.User
    if err := c.ShouldBindJSON(&loginUser); err != nil {
        result := errors.HandleError(
            func() (interface{}, error) {
                return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
            },
            "binding JSON for login user",
        )
        result.RespondWithJSON(c)
        return
    }

    // 2. Call business logic layer (Service)
    result := errors.HandleError(
        func() (interface{}, error) {
            user, err := h.authService.Login(loginUser.Email, loginUser.Password)
            if err != nil {
                return nil, err
            }
            return user, nil
        },
        "logging in user",
    )
    
    // 3. Send response
    result.RespondWithJSON(c)
}
```

**Handler Responsibilities**:
- Parse and validate HTTP request data
- Call appropriate service methods
- Handle errors using the error handling system
- Send HTTP responses

### 5. Service Layer (`internal/service/auth_service.go`)

The handler delegates business logic to the service layer:

```go
func (s *authService) Login(email, password string) (*models.User, error) {
    // 1. Fetch user from database via repository
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid credentials", nil)
    }

    // 2. Verify password using bcrypt
    passwordIsValid := utils.CheckPasswordHash(user.Password, password)
    if !passwordIsValid {
        return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid credentials", nil)
    }

    // 3. Clear password from response
    user.Password = ""
    
    // 4. Generate JWT tokens
    accessToken, err := utils.GenerateJwtToken(user.Email, false)
    if err != nil {
        return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token", err)
    }
    
    refreshToken, err := utils.GenerateJwtToken(user.Email, true)
    if err != nil {
        return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token", err)
    }
    
    // 5. Attach tokens to user object
    user.AccessToken = accessToken
    user.RefreshToken = refreshToken
    
    return user, nil
}
```

**Service Responsibilities**:
- Implement business logic and rules
- Coordinate between different repositories
- Handle authentication and authorization
- Generate tokens and perform security operations

### 6. Repository Layer (`internal/repository/user_repo.go`)

The service calls the repository for data access:

```go
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    
    // GORM query to find user by email
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "User not found", err)
        }
        return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user by email", err)
    }
    
    return &user, nil
}
```

**Repository Responsibilities**:
- Handle all database operations
- Convert database errors to application errors
- Implement data access patterns
- Use GORM for ORM operations

### 7. Database Layer (GORM + SQLite)

GORM translates the repository call into SQL:

```sql
-- Generated SQL query
SELECT * FROM users WHERE email = 'user@example.com' LIMIT 1;
```

The database returns the user record, which GORM maps back to the `User` struct.

### 8. Utility Functions

#### Password Verification (`internal/utils/hash.go`)
```go
func CheckPasswordHash(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

#### JWT Token Generation (`internal/utils/jwt_generator.go`)
```go
func GenerateJwtToken(email string, isRefresh bool) (string, error) {
    var expirationTime time.Time
    if isRefresh {
        expirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 days
    } else {
        expirationTime = time.Now().Add(15 * time.Minute)   // 15 minutes
    }
    
    claims := &jwt.MapClaims{
        "sub": email,
        "exp": expirationTime.Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}
```

### 9. Error Handling System (`pkg/errors/`)

The application uses a sophisticated error handling system:

#### Custom HTTP Error (`pkg/errors/errors.go`)
```go
type HTTPError struct {
    StatusCode int
    Message    string
    Err        error
}

func NewHTTPError(statusCode int, message string, err error) *HTTPError {
    return &HTTPError{
        StatusCode: statusCode,
        Message:    message,
        Err:        err,
    }
}
```

#### Error Handler (`pkg/errors/handler.go`)
```go
func HandleError(operation func() (interface{}, error), operationName string) *OperationResult {
    result := &OperationResult{
        OperationName: operationName,
        SuccessCode:   http.StatusOK,
        ErrorCode:     http.StatusInternalServerError,
    }

    data, err := operation()
    if err != nil {
        // Extract status code and message from HTTPError
        if statusCode, hasCode := GetStatusCode(err); hasCode {
            result.ErrorCode = statusCode
        }
        if message, hasMessage := GetErrorMessage(err); hasMessage {
            result.ErrorMessage = message
        }
        result.Error = err
    } else {
        result.Data = data
        result.SuccessMessage = fmt.Sprintf("%s completed successfully", operationName)
    }
    return result
}
```

#### Response Builder
```go
func (r *OperationResult) RespondWithJSON(c *gin.Context) {
    if r.Error != nil {
        c.JSON(r.ErrorCode, gin.H{
            "success":        false,
            "error":          r.ErrorMessage,
            "verbose_error":  r.VerboseErrorMessage,
            "status_code":    r.ErrorCode,
        })
    } else {
        response := gin.H{
            "success":     true,
            "message":     r.SuccessMessage,
            "status_code": r.SuccessCode,
        }
        if r.Data != nil {
            response["data"] = r.Data
        }
        c.JSON(r.SuccessCode, response)
    }
}
```

### 10. HTTP Response

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "logging in user completed successfully",
  "status_code": 200,
  "data": {
    "id": "user-uuid-123",
    "firstName": "John",
    "lastName": "Doe",
    "email": "user@example.com",
    "phoneNumber": "+1234567890",
    "emailVerified": true,
    "firstTimeLogin": false,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### Error Response (401 Unauthorized)
```json
{
  "success": false,
  "error": "Invalid credentials",
  "verbose_error": "Detailed error information: Invalid credentials",
  "status_code": 401
}
```

## Request Flow Summary

1. **HTTP Request** → Gin Router
2. **Middleware Chain** → Logging, CORS, Recovery, Rate Limiting
3. **Route Matching** → `/api/v1/auth/login` → `authHandler.Login`
4. **Handler** → Parse JSON, validate input
5. **Service** → Business logic, password verification, token generation
6. **Repository** → Database query via GORM
7. **Database** → SQLite query execution
8. **Response Chain** → Repository → Service → Handler → Error Handler → HTTP Response

## Key Design Patterns

### 1. Dependency Injection
- Services are injected into handlers
- Repositories are injected into services
- Promotes testability and loose coupling

### 2. Repository Pattern
- Abstracts database operations
- Provides clean interface for data access
- Easy to mock for testing

### 3. Service Layer Pattern
- Encapsulates business logic
- Coordinates between repositories
- Handles complex operations

### 4. Middleware Pattern
- Cross-cutting concerns (logging, CORS, rate limiting)
- Request/response processing pipeline
- Reusable components

### 5. Error Handling Pattern
- Consistent error responses
- HTTP status code mapping
- Detailed error information for debugging

## Security Measures

1. **Password Hashing**: bcrypt with salt
2. **JWT Tokens**: Signed tokens with expiration
3. **Rate Limiting**: IP-based request limiting
4. **Input Validation**: JSON binding with validation
5. **Error Masking**: Generic error messages for security

This architecture provides a robust, scalable, and maintainable foundation for the dfood API, with clear separation of concerns and comprehensive error handling.