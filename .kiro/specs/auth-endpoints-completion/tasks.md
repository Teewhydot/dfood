# Implementation Plan

- [ ] 1. Create new request and response models
  - Add `PasswordResetRequest` struct to `internal/models/requests.go`
  - Add `UserResponse` struct to `internal/models/requests.go` for sanitized user data
  - Write unit tests for model validation
  - _Requirements: 3.1, 4.4_

- [ ] 2. Extend repository interface for account deletion
  - Add `Delete(id string) error` method to `UserRepository` interface in `internal/repository/interfaces.go`
  - Add `DeleteUserData(userID string) error` method for cascading deletes
  - _Requirements: 2.4_

- [ ] 3. Implement repository methods for user deletion
  - Create implementation of `Delete` method in user repository
  - Create implementation of `DeleteUserData` method to handle cascading deletes (addresses, permissions, etc.)
  - Write unit tests for deletion methods
  - _Requirements: 2.4, 2.5_

- [ ] 4. Extend AuthService interface with new methods
  - Add `SendPasswordReset(email string) error` method to `AuthService` interface
  - Add `GetCurrentUser(token string) (*models.User, error)` method to `AuthService` interface
  - _Requirements: 3.1, 4.1_

- [ ] 5. Implement GetCurrentUser service method
  - Create `GetCurrentUser` method in `authService` struct
  - Implement token validation and user data retrieval
  - Sanitize user data by removing sensitive fields
  - Write unit tests for user retrieval and data sanitization
  - _Requirements: 4.1, 4.3, 4.4, 4.5_

- [ ] 6. Implement SendPasswordReset service method
  - Create `SendPasswordReset` method in `authService` struct
  - Generate secure reset token using JWT system
  - Integrate with existing email service to send reset emails
  - Implement email enumeration prevention (return success regardless of email existence)
  - Write unit tests for password reset flow
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ] 7. Implement GetCurrentUser handler
  - Update `GetCurrentUser` method in `AuthHandler` to extract JWT token from request
  - Call service method and handle response
  - Implement proper error handling using existing error patterns
  - Write unit tests for handler with valid/invalid tokens
  - _Requirements: 4.1, 4.2, 4.5_

- [ ] 8. Implement Logout handler
  - Update `Logout` method in `AuthHandler` to extract JWT token from request
  - Call existing service method for token invalidation
  - Implement proper error handling and success response
  - Write unit tests for logout functionality
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [ ] 9. Implement SendPasswordReset handler
  - Update `SendPasswordReset` method in `AuthHandler` to bind JSON request
  - Validate email format using existing validation patterns
  - Call service method and return generic success response
  - Write unit tests for password reset endpoint
  - _Requirements: 3.1, 3.2, 3.5_

- [ ] 10. Implement DeleteAccount handler
  - Update `DeleteAccount` method in `AuthHandler` to extract JWT token
  - Extract user email from token claims
  - Call existing service method for account deletion
  - Implement proper error handling and confirmation response
  - Write unit tests for account deletion endpoint
  - _Requirements: 2.1, 2.2, 2.3, 2.5_

- [ ] 11. Add integration tests for complete authentication flows
  - Write tests for login → use protected endpoint → logout flow
  - Write tests for account deletion and data cleanup verification
  - Write tests for password reset email generation
  - Write tests for token invalidation after logout
  - _Requirements: 1.3, 1.4, 2.4, 3.4_

- [ ] 12. Add middleware integration and route testing
  - Verify all endpoints work with existing authentication middleware
  - Test error responses for invalid/expired tokens
  - Test successful responses for valid tokens
  - Ensure consistent error handling across all endpoints
  - _Requirements: 1.2, 2.2, 4.2, 4.5_