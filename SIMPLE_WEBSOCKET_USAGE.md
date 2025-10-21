# Simple WebSocket Usage Guide

## 🎯 How to Use the Simplified WebSocket System

This guide shows you exactly how to use the simplified WebSocket implementation.

---

## 1. Connect from Flutter Client

```dart
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

class SimpleWebSocketService {
  WebSocketChannel? _channel;
  
  // Connect to user profile updates
  Future<void> connectToUserUpdates(String userId) async {
    try {
      // Connect to the WebSocket endpoint
      _channel = WebSocketChannel.connect(
        Uri.parse('ws://localhost:8080/api/v1/users/$userId/watch')
      );
      
      // Listen for messages
      _channel!.stream.listen(
        (message) {
          final data = jsonDecode(message);
          _handleMessage(data);
        },
        onError: (error) {
          print('WebSocket error: $error');
        },
        onDone: () {
          print('WebSocket connection closed');
        },
      );
      
      print('Connected to WebSocket for user: $userId');
    } catch (e) {
      print('Failed to connect: $e');
    }
  }
  
  // Handle incoming messages
  void _handleMessage(Map<String, dynamic> message) {
    final type = message['type'];
    
    switch (type) {
      case 'connected':
        print('✅ WebSocket connected successfully');
        break;
        
      case 'user_update':
        final data = message['data'];
        final user = data['user'];
        final changes = data['changes'];
        
        print('👤 User updated: ${user['first_name']} ${user['last_name']}');
        print('📝 Changes: $changes');
        
        // Update your UI here
        _updateUserInUI(user, changes);
        break;
        
      case 'error':
        print('❌ WebSocket error: ${message['data']}');
        break;
        
      default:
        print('Unknown message type: $type');
    }
  }
  
  // Update UI with new user data
  void _updateUserInUI(Map<String, dynamic> user, Map<String, dynamic> changes) {
    // Example: Update user provider or state management
    // UserProvider.instance.updateUser(User.fromJson(user));
    
    // Show notification about what changed
    changes.forEach((field, value) {
      print('Field "$field" changed to: $value');
    });
  }
  
  // Disconnect
  void disconnect() {
    _channel?.sink.close();
    _channel = null;
  }
}
```

---

## 2. Test the WebSocket Connection

### Step 1: Start your Go server
```bash
go run cmd/main.go
```

### Step 2: Connect to WebSocket (using wscat for testing)
```bash
# Install wscat if you don't have it
npm install -g wscat

# Connect to user profile WebSocket
wscat -c ws://localhost:8080/api/v1/users/user123/watch
```

**Expected response:**
```json
{
  "type": "connected",
  "data": {
    "message": "Connected successfully"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Step 3: Update user profile via REST API
```bash
# Update user's first name
curl -X PUT http://localhost:8080/api/v1/users/user123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "first_name": "John Updated",
    "bio": "New bio text"
  }'
```

**WebSocket receives:**
```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user123",
      "first_name": "John Updated",
      "last_name": "Doe",
      "email": "john@example.com",
      "bio": "New bio text"
    },
    "changes": {
      "first_name": "John Updated",
      "bio": "New bio text"
    }
  },
  "timestamp": "2024-01-15T10:31:00Z"
}
```

### Step 4: Update single field
```bash
# Update just the phone number
curl -X PATCH http://localhost:8080/api/v1/users/user123/phone_number \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "value": "+1234567890"
  }'
```

**WebSocket receives:**
```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user123",
      "first_name": "John Updated",
      "phone_number": "+1234567890"
    },
    "changes": {
      "phone_number": "+1234567890"
    }
  },
  "timestamp": "2024-01-15T10:32:00Z"
}
```

---

## 3. Check Connection Statistics

```bash
# See how many users are connected
curl -X GET http://localhost:8080/api/v1/websocket/stats
```

**Response:**
```json
{
  "success": true,
  "connected_users": 1,
  "message": "WebSocket connection statistics"
}
```

---

## 4. Flutter Integration Example

```dart
class UserProfileScreen extends StatefulWidget {
  final String userId;
  
  const UserProfileScreen({Key? key, required this.userId}) : super(key: key);
  
  @override
  State<UserProfileScreen> createState() => _UserProfileScreenState();
}

class _UserProfileScreenState extends State<UserProfileScreen> {
  final SimpleWebSocketService _wsService = SimpleWebSocketService();
  User? _user;
  
  @override
  void initState() {
    super.initState();
    _loadUser();
    _connectWebSocket();
  }
  
  Future<void> _loadUser() async {
    // Load user data from API
    final user = await UserService.getUser(widget.userId);
    setState(() {
      _user = user;
    });
  }
  
  Future<void> _connectWebSocket() async {
    await _wsService.connectToUserUpdates(widget.userId);
    
    // Listen for user updates
    _wsService.userUpdateStream.listen((updatedUser) {
      setState(() {
        _user = updatedUser;
      });
      
      // Show snackbar about the update
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Profile updated!')),
      );
    });
  }
  
  @override
  Widget build(BuildContext context) {
    if (_user == null) {
      return const Center(child: CircularProgressIndicator());
    }
    
    return Scaffold(
      appBar: AppBar(
        title: Text('${_user!.firstName} ${_user!.lastName}'),
        actions: [
          // WebSocket connection indicator
          Icon(
            Icons.wifi,
            color: _wsService.isConnected ? Colors.green : Colors.red,
          ),
        ],
      ),
      body: Column(
        children: [
          ListTile(
            title: Text('Name'),
            subtitle: Text('${_user!.firstName} ${_user!.lastName}'),
          ),
          ListTile(
            title: Text('Email'),
            subtitle: Text(_user!.email),
          ),
          ListTile(
            title: Text('Phone'),
            subtitle: Text(_user!.phoneNumber ?? 'Not set'),
          ),
          if (_user!.bio != null)
            ListTile(
              title: Text('Bio'),
              subtitle: Text(_user!.bio!),
            ),
        ],
      ),
    );
  }
  
  @override
  void dispose() {
    _wsService.disconnect();
    super.dispose();
  }
}
```

---

## 5. What Happens Behind the Scenes

```
1. User opens profile screen in Flutter app
   ↓
2. Flutter connects to WebSocket: /api/v1/users/123/watch
   ↓
3. Go server validates user exists and upgrades connection
   ↓
4. Server sends "connected" message to Flutter
   ↓
5. User updates profile via another device/tab
   ↓
6. Go server saves changes to database
   ↓
7. Go server automatically sends "user_update" message to WebSocket
   ↓
8. Flutter receives message and updates UI instantly
   ↓
9. User sees changes without refreshing!
```

---

## 6. Key Benefits of This Simple Approach

✅ **Easy to understand** - Just 3 message types  
✅ **Easy to test** - One endpoint to connect to  
✅ **Easy to debug** - Simple connection mapping  
✅ **Production ready** - Handles errors and cleanup  
✅ **Flutter friendly** - Follows patterns you know  
✅ **Automatic updates** - No manual broadcasting needed  

---

## 7. Troubleshooting

### WebSocket won't connect
- Check if server is running on correct port
- Verify user ID exists in database
- Check for CORS issues (upgrader allows all origins in development)

### Not receiving updates
- Verify WebSocket connection is active
- Check server logs for errors
- Ensure user service has WebSocket service set

### Connection drops frequently
- Check network stability
- Implement reconnection logic in Flutter
- Monitor server logs for connection errors

This simplified implementation gives you real-time user profile updates with minimal complexity!