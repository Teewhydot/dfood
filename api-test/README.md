# API Testing Files

This directory contains HTTP test files for all implemented endpoints in the dfood API.

## Files Overview

- **`auth.http`** - Authentication endpoints (register, login, logout, password management)
- **`users.http`** - User profile management endpoints
- **`addresses.http`** - Address management endpoints
- **`restaurants.http`** - Restaurant discovery and search endpoints
- **`foods.http`** - Food/menu browsing and search endpoints
- **`orders.http`** - Order creation and management endpoints
- **`favorites.http`** - Favorites management endpoints
- **`notifications.http`** - Notification management endpoints
- **`payments.http`** - Payment endpoints (not implemented - external service)
- **`chats.http`** - Chat/messaging endpoints (not implemented - WebSocket)
- **`uploads.http`** - File upload endpoints (not implemented - file storage)
- **`workflow.http`** - Complete user journey workflow example

## How to Use

### Prerequisites
1. Start the dfood server: `go run cmd/main.go`
2. Server should be running on `http://localhost:8080`

### Using VS Code REST Client Extension
1. Install the "REST Client" extension in VS Code
2. Open any `.http` file
3. Click "Send Request" above each HTTP request
4. View responses in the right panel

### Authentication Flow
1. **Register**: Use `auth.http` to register a new user
2. **Login**: Use the login endpoint to get an access token
3. **Copy Token**: Copy the `access_token` from the login response
4. **Replace Token**: Replace `{{access_token}}` or `YOUR_ACCESS_TOKEN_HERE` in other requests

### Testing Workflow
1. Start with `workflow.http` for a complete user journey
2. Use individual files to test specific features
3. Update IDs in requests based on actual data from responses

## Implementation Status

### ✅ Fully Implemented
- Authentication (register, login, password update)
- User profile management
- Address management
- Restaurant discovery
- Food browsing
- Order management
- Favorites management
- Notification management

### ❌ Not Implemented (External Services)
- Email services (password reset, email verification)
- Push notifications (FCM)
- File uploads (image storage)
- Real-time features (WebSocket endpoints)
- Payment processing (external gateway)

## Sample Data

The test files use sample IDs like:
- `user-123`, `test-user-001`
- `restaurant-123`
- `food-123`
- `order-123`

Replace these with actual IDs from your database or API responses.

## Error Handling

Common HTTP status codes you might see:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `404` - Not Found
- `409` - Conflict (duplicate data)
- `500` - Internal Server Error
- `501` - Not Implemented

## Tips

1. **Start with workflow.http** for a complete test sequence
2. **Copy actual IDs** from responses to use in subsequent requests
3. **Check server logs** for detailed error information
4. **Use proper JSON formatting** in request bodies
5. **Include Authorization header** for protected endpoints