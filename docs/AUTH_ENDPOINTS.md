# Authentication Endpoints Documentation

## Overview
Complete authentication system with user registration, login, logout, password management, and account deletion.

## Endpoints

### 1. User Registration
**POST** `/api/v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "+1234567890"
}
```

**Response (201):**
```json
{
  "id": "user_id",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "+1234567890",
  "createdAt": "2025-10-05T14:30:00Z"
}
```

**Features:**
- Sends welcome email automatically
- Password hashing with bcrypt
- Duplicate email validation

### 2. User Login
**POST** `/api/v1/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "id": "user_id",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "accessToken": "jwt_access_token",
  "refreshToken": "jwt_refresh_token"
}
```

### 3. Get Current User
**GET** `/api/v1/auth/me`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "id": "user_id",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "+1234567890",
  "emailVerified": false,
  "firstTimeLogin": true
}
```

### 4. Update Password
**PUT** `/api/v1/auth/password`

**Request Body:**
```json
{
  "email": "user@example.com",
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response (200):**
```json
{
  "message": "Password updated successfully"
}
```

**Features:**
- Validates current password
- Ensures new password is different
- Secure password hashing

### 5. Forgot Password
**POST** `/api/v1/auth/forgot-password`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response (200):**
```json
{
  "message": "If an account with that email exists, a password reset link has been sent",
  "email": "user@example.com"
}
```

**Note:** Currently returns success message. Full implementation would generate reset token and send email.

### 6. Logout
**POST** `/api/v1/auth/logout`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "message": "Successfully logged out"
}
```

**Features:**
- Invalidates the provided JWT token
- Adds token to blacklist

### 7. Delete Account
**DELETE** `/api/v1/auth/account`

**Request Body:**
```json
{
  "email": "user@example.com",
  "token": "jwt_access_token"
}
```

**Response (200):**
```json
{
  "message": "Account successfully deleted"
}
```

**Note:** Currently returns success but doesn't actually delete from database. Implementation needed.

## Security Features

### JWT Token Management
- **Access Tokens**: 15-minute expiration
- **Refresh Tokens**: 7-day expiration
- **Token Blacklisting**: Invalidated tokens are blacklisted
- **Token Validation**: All protected endpoints validate tokens

### Password Security
- **Bcrypt Hashing**: Secure password storage
- **Password Validation**: Minimum length requirements
- **Current Password Check**: Required for password updates

### Email Integration
- **Welcome Emails**: Automatic on registration
- **SendGrid Integration**: Professional email delivery
- **Async Processing**: Non-blocking email sending

## Error Handling

All endpoints use consistent error handling:

**400 Bad Request:**
```json
{
  "error": "Invalid JSON payload",
  "details": "Validation error details"
}
```

**401 Unauthorized:**
```json
{
  "error": "Invalid credentials",
  "details": "Authentication failed"
}
```

**409 Conflict:**
```json
{
  "error": "User already exists",
  "details": "Email already registered"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error",
  "details": "Server processing error"
}
```

## Testing

Use the provided test file `api-test/auth-extended.http` to test all endpoints.

## Next Steps

1. **Complete Delete Account**: Implement actual user deletion from database
2. **Password Reset Flow**: Complete implementation with reset tokens
3. **Email Verification**: Add email verification on registration
4. **Rate Limiting**: Add specific rate limits for auth endpoints
5. **Audit Logging**: Log authentication events for security