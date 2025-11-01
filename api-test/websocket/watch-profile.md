# WebSocket: Watch User Profile Updates

## Overview
This WebSocket endpoint allows clients to watch for real-time profile updates for a specific user.

**Endpoint:** `ws://localhost:8080/api/v1/users/:userId/watch`

## How It Works

1. Client connects to WebSocket endpoint with a user ID
2. Server immediately sends the current user profile
3. When profile is updated (via PUT or PATCH), all connected clients receive the update
4. Connection stays open until client disconnects

## Testing with the HTML Client

Use the provided HTML test client:

```bash
# Open the test client in your browser
open websocket-test-client.html
```

## Testing Manually

### 1. Connect to WebSocket

**Using JavaScript:**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/users/user-123/watch');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
};
```

**Using curl (for testing connection):**
```bash
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
  http://localhost:8080/api/v1/users/user-123/watch
```

### 2. Update Profile (Trigger WebSocket Update)

In another terminal, update the user profile:

```bash
# Update entire profile
curl -X PUT http://localhost:8080/api/v1/users/user-123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "first_name": "John Updated",
    "last_name": "Doe Updated",
    "bio": "New bio"
  }'

# Or update a single field
curl -X PATCH http://localhost:8080/api/v1/users/user-123/first_name \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"value": "Johnny"}'
```

### 3. Observe WebSocket Message

The WebSocket client will receive a message like:

```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user-123",
      "first_name": "Johnny",
      "last_name": "Doe Updated",
      "email": "john@example.com",
      "bio": "New bio",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "changes": {
      "first_name": "Johnny"
    }
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Message Types

### Initial Connection
When you first connect, you receive the current profile:
```json
{
  "type": "user_update",
  "data": {
    "user": { ... },
    "changes": { "initial": true }
  },
  "timestamp": "..."
}
```

### Profile Update
When profile is updated:
```json
{
  "type": "user_update",
  "data": {
    "user": { ... },
    "changes": { "first_name": "New Name" }
  },
  "timestamp": "..."
}
```

### Connection Confirmed
Server acknowledges connection:
```json
{
  "type": "connected",
  "data": {
    "message": "Connected successfully"
  },
  "timestamp": "..."
}
```

## Testing Tools

### Postman
1. Create new WebSocket Request
2. Set URL: `ws://localhost:8080/api/v1/users/user-123/watch`
3. Click Connect
4. Update profile via REST API in another tab
5. Watch messages arrive in WebSocket tab

### Browser DevTools
1. Open browser console (F12)
2. Paste the JavaScript code above
3. Update profile via REST API
4. Watch console for messages

### wscat (CLI tool)
```bash
npm install -g wscat
wscat -c ws://localhost:8080/api/v1/users/user-123/watch
```

## Error Cases

**User not found:**
```json
{
  "success": false,
  "error": "User not found"
}
```
Connection is rejected with HTTP 404.

**Invalid user ID:**
```json
{
  "success": false,
  "error": "User ID is required"
}
```
Connection is rejected with HTTP 400.

## Stats Endpoint

Check WebSocket connection statistics:

```bash
GET http://localhost:8080/api/v1/websocket/stats
```

Response:
```json
{
  "success": true,
  "data": {
    "total_connections": 3,
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```
