# WebSocket Quick Start Guide

## What We Built

A complete WebSocket endpoint that broadcasts user profile updates in real-time.

**Endpoint:** `ws://localhost:8080/api/v1/users/:userId/watch`

## Quick Test (3 Steps)

### Step 1: Start the Server

```bash
cd /mnt/c/Users/Hp/IdeaProjects/dfood
go run cmd/main.go
```

### Step 2: Open the Test Client

Open `websocket-test-client.html` in your browser.

### Step 3: Test the WebSocket

1. **Connect:**
   - Enter user ID (e.g., `user-123`)
   - Click "Connect"
   - You'll see the current profile

2. **Update Profile:**
   - In another terminal, first register and login to get a token:
   ```bash
   # Register
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "first_name": "John",
       "last_name": "Doe",
       "email": "john@example.com",
       "phone_number": "+1234567890",
       "password": "password123"
     }'

   # Login (save the access_token from response)
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "john@example.com",
       "password": "password123"
     }'
   ```

3. **Watch Updates in Real-Time:**
   ```bash
   # Use the token from login
   export TOKEN="your_access_token_here"
   export USER_ID="user_id_from_register_response"

   # Update profile - watch WebSocket client receive update!
   curl -X PUT http://localhost:8080/api/v1/users/$USER_ID \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{
       "first_name": "Jane",
       "last_name": "Smith",
       "bio": "Updated via REST API!"
     }'
   ```

4. **See the Magic:**
   - The WebSocket client instantly shows the profile update!
   - No polling, no refreshing - pure real-time updates

## How It Works

```
┌─────────┐                    ┌────────┐                    ┌──────────┐
│ Client  │◄───WebSocket───────┤ Server │◄────REST API───────┤  Client  │
│   #1    │    (watching)      │        │   (updating)       │    #2    │
└─────────┘                    └────────┘                    └──────────┘
     │                              │                              │
     │  1. Connect & get profile    │                              │
     │◄─────────────────────────────│                              │
     │                              │                              │
     │                              │  2. Update profile via PUT   │
     │                              │◄─────────────────────────────│
     │                              │                              │
     │  3. Instant update broadcast │  4. Success response         │
     │◄─────────────────────────────│─────────────────────────────►│
     │     (real-time!)             │                              │
```

## Message Format

### When You Connect:
```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user-123",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com"
    },
    "changes": { "initial": true }
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### When Profile Updates:
```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user-123",
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "john@example.com"
    },
    "changes": {
      "first_name": "Jane",
      "last_name": "Smith"
    }
  },
  "timestamp": "2024-01-01T12:05:00Z"
}
```

## Architecture

### Components Created:

1. **WebSocket Service** (`internal/service/websocket.go`)
   - Manages all WebSocket connections
   - Broadcasts updates to connected clients
   - One connection per user (single device)

2. **WebSocket Handler** (`internal/api/handlers/websocket.go`)
   - Upgrades HTTP to WebSocket
   - Sends initial profile data
   - Handles connection lifecycle

3. **User Service Integration** (`internal/service/user_service.go`)
   - Broadcasts on `Update()` (line 96-102)
   - Broadcasts on `UpdateField()` (line 140-148)

4. **Routes** (`internal/api/routes/routes.go`)
   - `GET /api/v1/users/:userId/watch` - WebSocket endpoint
   - `GET /api/v1/websocket/stats` - Connection stats

5. **Test Client** (`api-test/websocket/websocket-test-client.html`)
   - Beautiful web UI for testing
   - Real-time message display
   - Profile update simulator

## Multiple Clients

Try this:
1. Open test client in multiple browser tabs
2. Connect all tabs to same user ID
3. Update profile via REST API
4. Watch ALL tabs receive the update simultaneously!

## Production Considerations

Before deploying:

1. **Security:**
   - Add authentication to WebSocket endpoint
   - Validate user can only watch their own profile
   - Implement rate limiting

2. **Scaling:**
   - Current implementation: single server, in-memory
   - For multiple servers: use Redis Pub/Sub
   - Consider message queues for reliability

3. **Connection Limits:**
   - Add max connections per user
   - Implement connection timeouts
   - Add ping/pong heartbeat

## Troubleshooting

**Connection refused?**
- Make sure server is running
- Check port 8080 is not blocked

**Not receiving updates?**
- Verify user ID exists
- Check auth token is valid
- Look at server logs

**CORS issues?**
- Already configured to allow all origins
- For production, configure properly in routes

## Next Steps

Want to implement more WebSocket endpoints?
- Order tracking (real-time delivery status)
- Notifications (live notification feed)
- Chat (real-time messaging)

The pattern is the same:
1. Use `WebSocketService` to manage connections
2. Broadcast from service layer when data changes
3. Create handler to upgrade connection
4. Add route

Perfect template for any real-time feature!
