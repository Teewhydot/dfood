# Design Document

## Overview

This design document outlines the implementation of four remaining authentication endpoints in the dfood application: Logout, DeleteAccount, SendPasswordReset, and GetCurrentUser. The implementation will follow the existing clean architecture pattern with handlers, services, and repositories, leveraging the current JWT token management system and email service infrastructure.

## Architecture

The implementation follows the existing layered architecture:

```
Handler Layer (API) → Service Layer (Business Logic) → Repository Layer (Data Access)
```

### Current System Integration Points

- **JWT Token Management**: Utilizes existing `utils/jwt_generator.go` with token blacklisting
- **Email Service**: Leverages existing `EmailService` interface with SendGrid integration
- **Error Handling**: Uses existing `pkg/errors` package for consistent error responses
- **User Repository**: Extends existing `UserRepository` interface for account deletion
- **Authentication Middleware**: Works with existing JWT validation middleware

## Components and Interfaces

### 1. Handler Layer Updates

**File**: `internal/api/handlers/auth.go`

The AuthHandler will be updated to implement the four TODO methods:

```go
func (h *AuthHandler) Logout(c *gin.Context)
func (h *AuthHandler) DeleteAccount(c *gin.Context) 
func (h *AuthHandler) SendPasswordReset(c *gin.Context)
func (h *AuthHandler) GetCurrentUser(c *gin.Context)
```

### 2. Service Layer Updates

**File**: `internal/service/auth_service.go`

The AuthService interface already contains `Logout` and `DeleteAccount` methods. Additional methods will be added:

```go
type AuthService interface {
    // Existing methods...
    Logout(token string) error
    DeleteAccount(email, token string) error
    
    // New methods to be added:
    SendPasswordReset(email string) error
    GetCurrentUser(token string) (*models.User, error)
}
```

### 3. Repository Layer Updates

**File**: `internal/repository/interfaces.go`

The UserRepository interface will be extended to support account deletion:

```go
type UserRepository interface {
    // Existing methods...
    
    // New method for account deletion:
    Delete(id string) error
    DeleteUserData(userID string) error // For cascading deletes
}
```

### 4. New Request/Response Models

**File**: `internal/models/requests.go`

```go
// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
    Email string `json:"email" binding:"required,email"`
}

// UserResponse represents sanitized user data for API responses
type UserResponse struct {
    ID              string    `json:"id"`
    FirstName       string    `json:"first_name"`
    LastName        string    `json:"last_name"`
    Email           string    `json:"email"`
    PhoneNumber     string    `json:"phone_number"`
    ProfileImageURL *string   `json:"profile_image_url,omitempty"`
    Bio             *string   `json:"bio,omitempty"`
    FirstTimeLogin  bool      `json:"first_time_login"`
    EmailVerified   bool      `json:"email_verified"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

## Data Models

### Token Management

The existing JWT token system will be leveraged:
- **Access Tokens**: 15-minute expiration
- **Refresh Tokens**: 7-day expiration  
- **Token Blacklisting**: In-memory blacklist (production should use Redis)
- **Token Validation**: Existing `ValidateToken` function

### Password Reset Flow

```go
type PasswordResetToken struct {
    Email     string
    Token     string
    ExpiresAt time.Time
}
```

Note: For this implementation, password reset tokens will be generated using the existing JWT system with a special claim to identify them as reset tokens.

### User Data Sanitization

Sensitive fields will be excluded from API responses:
- Password (already handled)
- Internal database IDs where not needed
- FCM tokens
- Refresh tokens (except during login)

## Error Handling

### Consistent Error Responses

All endpoints will use the existing `errors.HandleError` pattern:

```go
result := errors.HandleError(
    func() (interface{}, error) {
        // Business logic
        return data, err
    },
    "operation description",
)
result.RespondWithJSON(c)
```

### Security Considerations

1. **Email Enumeration Prevention**: Password reset endpoint returns success regardless of email existence
2. **Token Validation**: All protected endpoints validate JWT tokens
3. **Rate Limiting**: Password reset requests should be rate-limited (implementation note for future)
4. **Secure Token Generation**: Use existing JWT system for reset tokens
5. **Data Sanitization**: Remove sensitive data from responses

## Testing Strategy

### Unit Tests

1. **Handler Tests**:
   - Test each endpoint with valid/invalid tokens
   - Test request validation and error handling
   - Mock service layer dependencies

2. **Service Tests**:
   - Test business logic for each operation
   - Test token validation and blacklisting
   - Test email sending functionality
   - Test user data retrieval and sanitization

3. **Integration Tests**:
   - Test complete authentication flows
   - Test token lifecycle (login → use → logout)
   - Test account deletion cascade effects

### Test Data Setup

- Create test users with known credentials
- Generate test JWT tokens for validation
- Mock email service for password reset testing
- Test database cleanup after account deletion

### Security Testing

- Test token invalidation after logout
- Test access denial with blacklisted tokens
- Test account deletion data cleanup
- Test password reset token expiration
- Test unauthorized access attempts

## Implementation Sequence

### Phase 1: Core Infrastructure
1. Add new request/response models
2. Extend repository interfaces
3. Update service interfaces

### Phase 2: Service Layer Implementation
1. Implement `GetCurrentUser` service method
2. Implement `SendPasswordReset` service method
3. Extend user repository with delete methods

### Phase 3: Handler Implementation
1. Implement `GetCurrentUser` handler
2. Implement `Logout` handler (already has service method)
3. Implement `SendPasswordReset` handler
4. Implement `DeleteAccount` handler (already has service method)

### Phase 4: Testing and Validation
1. Unit tests for all new methods
2. Integration testing
3. Security validation
4. Error handling verification

## API Endpoint Specifications

### GET /auth/me
- **Purpose**: Get current authenticated user
- **Authentication**: Required (JWT token)
- **Response**: Sanitized user profile data

### POST /auth/logout  
- **Purpose**: Logout current user
- **Authentication**: Required (JWT token)
- **Response**: Success confirmation

### DELETE /auth/account
- **Purpose**: Delete user account permanently
- **Authentication**: Required (JWT token)
- **Response**: Deletion confirmation

### POST /auth/password-reset
- **Purpose**: Send password reset email
- **Authentication**: Not required
- **Request Body**: Email address
- **Response**: Generic success message

## Dependencies

### Existing Dependencies
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `github.com/sendgrid/sendgrid-go` - Email service
- `gorm.io/gorm` - Database operations
- `github.com/gin-gonic/gin` - HTTP framework

### No New Dependencies Required
The implementation uses existing infrastructure and dependencies.