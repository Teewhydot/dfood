# Flutter WebSocket Implementation Guide

## 📚 Table of Contents
1. [Implemented WebSocket Endpoints](#implemented-websocket-endpoints)
2. [Message Formats](#message-formats)
3. [Flutter Setup](#flutter-setup)
4. [Complete Flutter Implementation](#complete-flutter-implementation)
5. [Usage Examples](#usage-examples)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

---

## Implemented WebSocket Endpoints

### Summary Table

| Endpoint | Purpose | Status | Auto-Reconnect |
|----------|---------|--------|----------------|
| `/api/v1/realtime/users/:userId/watch` | User profile updates | ✅ | ✅ |
| `/api/v1/realtime/users/:userId/addresses/watch` | Address CRUD | ✅ | ✅ |
| `/api/v1/realtime/users/:userId/favorites/foods/watch` | Favorite foods | ✅ | ✅ |
| `/api/v1/realtime/users/:userId/favorites/restaurants/watch` | Favorite restaurants | ✅ | ✅ |
| `/api/v1/realtime/users/:userId/notifications/watch` | Notifications | ✅ | ✅ |
| `/api/v1/realtime/stats` | Connection stats | ✅ | N/A |

### Endpoint Details

#### 1. User Profile Updates
```
Endpoint: ws://your-server.com/api/v1/realtime/users/{userId}/watch
Authentication: Required
Triggers: Profile updates, name changes, image uploads
```

#### 2. Address Updates
```
Endpoint: ws://your-server.com/api/v1/realtime/users/{userId}/addresses/watch
Authentication: Required
Triggers: Add, update, delete, set default address
```

#### 3. Favorite Foods Updates
```
Endpoint: ws://your-server.com/api/v1/realtime/users/{userId}/favorites/foods/watch
Authentication: Required
Triggers: Add/remove favorite foods, clear favorites
```

#### 4. Favorite Restaurants Updates
```
Endpoint: ws://your-server.com/api/v1/realtime/users/{userId}/favorites/restaurants/watch
Authentication: Required
Triggers: Add/remove favorite restaurants, clear favorites
```

#### 5. Notification Updates
```
Endpoint: ws://your-server.com/api/v1/realtime/users/{userId}/notifications/watch
Authentication: Required
Triggers: New notifications
```

---

## Message Formats

### Generic Message Structure

All messages follow this format:

```json
{
  "type": "message_type",
  "data": { /* specific data */ },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

### Message Types and Payloads

#### User Update Message

**Type:** `user_update`

```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user123",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "phone_number": "+1234567890",
      "bio": "Food lover",
      "profile_image_url": "https://..."
    },
    "changes": {
      "first_name": "John"
    }
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Address Add Message

**Type:** `address_add`

```json
{
  "type": "address_add",
  "data": {
    "address": {
      "id": "addr123",
      "user_id": "user123",
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip_code": "10001",
      "address": "123 Main St",
      "apartment": "Apt 5B",
      "title": "Home",
      "type": "home",
      "is_default": true
    },
    "action": "add"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Address Update Message

**Type:** `address_update`

```json
{
  "type": "address_update",
  "data": {
    "address": { /* full address object */ },
    "action": "update",
    "changes": {
      "apartment": "Apt 6C"
    }
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Address Delete Message

**Type:** `address_delete`

```json
{
  "type": "address_delete",
  "data": {
    "address": { /* address that was deleted */ },
    "action": "delete"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Address Default Message

**Type:** `address_default`

```json
{
  "type": "address_default",
  "data": {
    "address": { /* new default address */ },
    "action": "set_default"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Favorite Add Message

**Type:** `favorite_add`

```json
{
  "type": "favorite_add",
  "data": {
    "resource_type": "food",  // or "restaurant"
    "resource_id": "food123",
    "action": "add",
    "resource": {
      "id": "food123",
      "name": "Pizza",
      "price": 12.99,
      // ... full food/restaurant object
    }
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Favorite Remove Message

**Type:** `favorite_remove`

```json
{
  "type": "favorite_remove",
  "data": {
    "resource_type": "food",
    "resource_id": "food123",
    "action": "remove",
    "resource": { /* food/restaurant object */ }
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Favorites Clear Message

**Type:** `favorites_clear`

```json
{
  "type": "favorites_clear",
  "data": {
    "action": "clear"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Notification New Message

**Type:** `notification_new`

```json
{
  "type": "notification_new",
  "data": {
    "notification": {
      "id": "notif123",
      "user_id": "user123",
      "title": "Order Ready",
      "body": "Your order is ready for pickup",
      "type": "order",
      "is_read": false,
      "created_at": "2025-10-31T12:00:00Z"
    },
    "action": "new"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Connected Message

**Type:** `connected`

```json
{
  "type": "connected",
  "data": {
    "message": "Connected successfully"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

#### Error Message

**Type:** `error`

```json
{
  "type": "error",
  "data": {
    "message": "Error description",
    "code": "ERROR_CODE"
  },
  "timestamp": "2025-10-31T12:00:00Z"
}
```

---

## Flutter Setup

### 1. Add Dependencies

**File:** `pubspec.yaml`

```yaml
dependencies:
  flutter:
    sdk: flutter

  # WebSocket support
  web_socket_channel: ^2.4.0

  # State management (choose one)
  provider: ^6.1.1          # or
  riverpod: ^2.4.0          # or
  bloc: ^8.1.0              # your choice

  # HTTP client
  dio: ^5.3.2

  # Logging (optional but recommended)
  logger: ^2.0.0
```

### 2. File Structure

```
lib/
├── models/
│   ├── realtime_models.dart      # WebSocket message models
│   ├── user.dart                 # User model
│   ├── address.dart              # Address model
│   ├── food.dart                 # Food model
│   └── notification.dart         # Notification model
├── services/
│   ├── realtime_service.dart     # Main WebSocket service
│   └── api_service.dart          # HTTP API service
├── providers/
│   ├── user_provider.dart        # User state
│   ├── address_provider.dart     # Address state
│   └── favorites_provider.dart   # Favorites state
└── screens/
    ├── profile_screen.dart
    ├── address_screen.dart
    └── favorites_screen.dart
```

---

## Complete Flutter Implementation

### Step 1: Create Message Models

**File:** `lib/models/realtime_models.dart`

```dart
import 'dart:convert';

/// Message types for all WebSocket endpoints
class RealtimeMessageType {
  // User messages
  static const String userUpdate = 'user_update';
  static const String connected = 'connected';
  static const String error = 'error';

  // Address messages
  static const String addressAdd = 'address_add';
  static const String addressUpdate = 'address_update';
  static const String addressDelete = 'address_delete';
  static const String addressDefault = 'address_default';

  // Favorites messages
  static const String favoriteAdd = 'favorite_add';
  static const String favoriteRemove = 'favorite_remove';
  static const String favoritesClear = 'favorites_clear';

  // Notification messages
  static const String notificationNew = 'notification_new';
  static const String notificationRead = 'notification_read';
  static const String notificationDelete = 'notification_delete';
}

/// Main WebSocket message structure
class RealtimeMessage {
  final String type;
  final Map<String, dynamic>? data;
  final DateTime timestamp;

  RealtimeMessage({
    required this.type,
    this.data,
    required this.timestamp,
  });

  factory RealtimeMessage.fromJson(Map<String, dynamic> json) {
    return RealtimeMessage(
      type: json['type'] ?? '',
      data: json['data'] as Map<String, dynamic>?,
      timestamp: DateTime.parse(json['timestamp']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'type': type,
      'data': data,
      'timestamp': timestamp.toIso8601String(),
    };
  }
}

/// User update data
class UserUpdateData {
  final dynamic user;  // Use your User model
  final Map<String, dynamic> changes;

  UserUpdateData({
    required this.user,
    required this.changes,
  });

  factory UserUpdateData.fromJson(Map<String, dynamic> json) {
    return UserUpdateData(
      user: json['user'],
      changes: Map<String, dynamic>.from(json['changes']),
    );
  }
}

/// Address update data
class AddressUpdateData {
  final dynamic address;  // Use your Address model
  final String action;
  final Map<String, dynamic>? changes;

  AddressUpdateData({
    required this.address,
    required this.action,
    this.changes,
  });

  factory AddressUpdateData.fromJson(Map<String, dynamic> json) {
    return AddressUpdateData(
      address: json['address'],
      action: json['action'],
      changes: json['changes'] != null
          ? Map<String, dynamic>.from(json['changes'])
          : null,
    );
  }
}

/// Favorite update data
class FavoriteUpdateData {
  final String resourceType;
  final String resourceId;
  final String action;
  final dynamic resource;

  FavoriteUpdateData({
    required this.resourceType,
    required this.resourceId,
    required this.action,
    this.resource,
  });

  factory FavoriteUpdateData.fromJson(Map<String, dynamic> json) {
    return FavoriteUpdateData(
      resourceType: json['resource_type'],
      resourceId: json['resource_id'],
      action: json['action'],
      resource: json['resource'],
    );
  }
}

/// Notification update data
class NotificationUpdateData {
  final dynamic notification;  // Use your Notification model
  final String action;

  NotificationUpdateData({
    required this.notification,
    required this.action,
  });

  factory NotificationUpdateData.fromJson(Map<String, dynamic> json) {
    return NotificationUpdateData(
      notification: json['notification'],
      action: json['action'],
    );
  }
}
```

### Step 2: Create RealtimeService

**File:** `lib/services/realtime_service.dart`

```dart
import 'dart:async';
import 'dart:convert';
import 'dart:developer';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;
import '../models/realtime_models.dart';

/// Manages all WebSocket connections and real-time updates
class RealtimeService {
  // Connection management
  final Map<String, WebSocketChannel> _connections = {};
  final Map<String, StreamController<RealtimeMessage>> _messageControllers = {};
  final Map<String, Timer> _reconnectTimers = {};

  String? _baseUrl;
  String? _userId;
  String? _authToken;
  bool _isInitialized = false;

  // Stream controllers for different data types
  final _userUpdateController = StreamController<UserUpdateData>.broadcast();
  final _addressUpdateController = StreamController<AddressUpdateData>.broadcast();
  final _favoriteUpdateController = StreamController<FavoriteUpdateData>.broadcast();
  final _notificationUpdateController = StreamController<NotificationUpdateData>.broadcast();

  // Public streams
  Stream<UserUpdateData> get userUpdateStream => _userUpdateController.stream;
  Stream<AddressUpdateData> get addressUpdateStream => _addressUpdateController.stream;
  Stream<FavoriteUpdateData> get favoriteUpdateStream => _favoriteUpdateController.stream;
  Stream<NotificationUpdateData> get notificationUpdateStream => _notificationUpdateController.stream;

  /// Initialize the service with user credentials
  void initialize({
    required String userId,
    required String baseUrl,
    String? authToken,
  }) {
    _userId = userId;
    _baseUrl = baseUrl;
    _authToken = authToken;
    _isInitialized = true;

    log('RealtimeService initialized for user: $userId');
  }

  /// Connect to user profile updates
  Future<void> connectToUserUpdates() async {
    await _connectToEndpoint(
      'user',
      '/realtime/users/$_userId/watch',
    );
  }

  /// Connect to address updates
  Future<void> connectToAddressUpdates() async {
    await _connectToEndpoint(
      'address',
      '/realtime/users/$_userId/addresses/watch',
    );
  }

  /// Connect to favorite foods updates
  Future<void> connectToFavoriteFoodsUpdates() async {
    await _connectToEndpoint(
      'favorite_foods',
      '/realtime/users/$_userId/favorites/foods/watch',
    );
  }

  /// Connect to favorite restaurants updates
  Future<void> connectToFavoriteRestaurantsUpdates() async {
    await _connectToEndpoint(
      'favorite_restaurants',
      '/realtime/users/$_userId/favorites/restaurants/watch',
    );
  }

  /// Connect to notification updates
  Future<void> connectToNotificationUpdates() async {
    await _connectToEndpoint(
      'notifications',
      '/realtime/users/$_userId/notifications/watch',
    );
  }

  /// Connect to all endpoints at once
  Future<void> connectToAll() async {
    if (!_isInitialized) {
      log('RealtimeService: Not initialized. Call initialize() first.');
      return;
    }

    await Future.wait([
      connectToUserUpdates(),
      connectToAddressUpdates(),
      connectToFavoriteFoodsUpdates(),
      connectToFavoriteRestaurantsUpdates(),
      connectToNotificationUpdates(),
    ]);
  }

  /// Private method to connect to a specific endpoint
  Future<void> _connectToEndpoint(String key, String endpoint) async {
    if (!_isInitialized || _userId == null || _baseUrl == null) {
      log('RealtimeService: Not initialized');
      return;
    }

    try {
      // Convert HTTP URL to WebSocket URL
      final wsUrl = _baseUrl!
          .replaceFirst('http://', 'ws://')
          .replaceFirst('https://', 'wss://');

      // Build full WebSocket URL
      final uri = Uri.parse('$wsUrl/api/v1$endpoint');

      log('RealtimeService: Connecting to $uri');

      // Create WebSocket channel
      // Note: web_socket_channel doesn't support custom headers directly
      // For authentication, you may need to pass token as query parameter
      // or handle it differently based on your backend implementation
      final channel = WebSocketChannel.connect(uri);

      _connections[key] = channel;

      // Create message controller for this connection
      final messageController = StreamController<RealtimeMessage>.broadcast();
      _messageControllers[key] = messageController;

      // Listen to messages
      channel.stream.listen(
        (message) => _handleMessage(key, message),
        onError: (error) => _handleError(key, error),
        onDone: () => _handleDisconnection(key),
      );

      log('RealtimeService: Connected to $key endpoint');
    } catch (e) {
      log('RealtimeService: Connection failed for $key - $e');
      _attemptReconnect(key, endpoint);
    }
  }

  /// Handle incoming WebSocket messages
  void _handleMessage(String connectionKey, dynamic message) {
    try {
      final Map<String, dynamic> json = jsonDecode(message);
      final realtimeMessage = RealtimeMessage.fromJson(json);

      log('RealtimeService: Received ${realtimeMessage.type} from $connectionKey');

      // Add to message stream
      _messageControllers[connectionKey]?.add(realtimeMessage);

      // Handle specific message types
      switch (realtimeMessage.type) {
        case RealtimeMessageType.connected:
          log('RealtimeService: $connectionKey connected');
          break;

        case RealtimeMessageType.userUpdate:
          if (realtimeMessage.data != null) {
            final updateData = UserUpdateData.fromJson(realtimeMessage.data!);
            _userUpdateController.add(updateData);
          }
          break;

        case RealtimeMessageType.addressAdd:
        case RealtimeMessageType.addressUpdate:
        case RealtimeMessageType.addressDelete:
        case RealtimeMessageType.addressDefault:
          if (realtimeMessage.data != null) {
            final updateData = AddressUpdateData.fromJson(realtimeMessage.data!);
            _addressUpdateController.add(updateData);
          }
          break;

        case RealtimeMessageType.favoriteAdd:
        case RealtimeMessageType.favoriteRemove:
        case RealtimeMessageType.favoritesClear:
          if (realtimeMessage.data != null) {
            final updateData = FavoriteUpdateData.fromJson(realtimeMessage.data!);
            _favoriteUpdateController.add(updateData);
          }
          break;

        case RealtimeMessageType.notificationNew:
        case RealtimeMessageType.notificationRead:
        case RealtimeMessageType.notificationDelete:
          if (realtimeMessage.data != null) {
            final updateData = NotificationUpdateData.fromJson(realtimeMessage.data!);
            _notificationUpdateController.add(updateData);
          }
          break;

        default:
          log('RealtimeService: Unknown message type: ${realtimeMessage.type}');
      }
    } catch (e) {
      log('RealtimeService: Error parsing message from $connectionKey - $e');
    }
  }

  /// Handle WebSocket errors
  void _handleError(String connectionKey, dynamic error) {
    log('RealtimeService: Error on $connectionKey - $error');
    final endpoint = _getEndpointForKey(connectionKey);
    if (endpoint != null) {
      _attemptReconnect(connectionKey, endpoint);
    }
  }

  /// Handle WebSocket disconnection
  void _handleDisconnection(String connectionKey) {
    log('RealtimeService: $connectionKey disconnected');
    final endpoint = _getEndpointForKey(connectionKey);
    if (endpoint != null) {
      _attemptReconnect(connectionKey, endpoint);
    }
  }

  /// Attempt to reconnect after disconnect
  void _attemptReconnect(String connectionKey, String endpoint) {
    _reconnectTimers[connectionKey]?.cancel();
    _reconnectTimers[connectionKey] = Timer(const Duration(seconds: 3), () {
      log('RealtimeService: Attempting to reconnect $connectionKey');
      _connectToEndpoint(connectionKey, endpoint);
    });
  }

  /// Get endpoint path for a connection key
  String? _getEndpointForKey(String key) {
    switch (key) {
      case 'user':
        return '/realtime/users/$_userId/watch';
      case 'address':
        return '/realtime/users/$_userId/addresses/watch';
      case 'favorite_foods':
        return '/realtime/users/$_userId/favorites/foods/watch';
      case 'favorite_restaurants':
        return '/realtime/users/$_userId/favorites/restaurants/watch';
      case 'notifications':
        return '/realtime/users/$_userId/notifications/watch';
      default:
        return null;
    }
  }

  /// Check if a specific connection is active
  bool isConnected(String connectionKey) {
    return _connections[connectionKey] != null;
  }

  /// Check if any connection is active
  bool get isAnyConnected {
    return _connections.isNotEmpty;
  }

  /// Get number of active connections
  int get connectionCount {
    return _connections.length;
  }

  /// Disconnect a specific connection
  void disconnect(String connectionKey) {
    _connections[connectionKey]?.sink.close(status.goingAway);
    _connections.remove(connectionKey);
    _messageControllers[connectionKey]?.close();
    _messageControllers.remove(connectionKey);
    _reconnectTimers[connectionKey]?.cancel();
    _reconnectTimers.remove(connectionKey);
  }

  /// Disconnect all connections
  void disconnectAll() {
    final keys = List<String>.from(_connections.keys);
    for (final key in keys) {
      disconnect(key);
    }
  }

  /// Dispose the service and clean up resources
  void dispose() {
    disconnectAll();
    _userUpdateController.close();
    _addressUpdateController.close();
    _favoriteUpdateController.close();
    _notificationUpdateController.close();
  }
}
```

### Step 3: Create State Provider (using Provider package)

**File:** `lib/providers/realtime_provider.dart`

```dart
import 'package:flutter/foundation.dart';
import '../services/realtime_service.dart';
import '../models/realtime_models.dart';

class RealtimeProvider with ChangeNotifier {
  final RealtimeService _realtimeService;

  RealtimeProvider(this._realtimeService) {
    _setupListeners();
  }

  // Connection status
  bool get isConnected => _realtimeService.isAnyConnected;
  int get connectionCount => _realtimeService.connectionCount;

  // Latest updates (optional - store if you want to show them)
  UserUpdateData? _latestUserUpdate;
  AddressUpdateData? _latestAddressUpdate;
  FavoriteUpdateData? _latestFavoriteUpdate;
  NotificationUpdateData? _latestNotificationUpdate;

  UserUpdateData? get latestUserUpdate => _latestUserUpdate;
  AddressUpdateData? get latestAddressUpdate => _latestAddressUpdate;
  FavoriteUpdateData? get latestFavoriteUpdate => _latestFavoriteUpdate;
  NotificationUpdateData? get latestNotificationUpdate => _latestNotificationUpdate;

  void _setupListeners() {
    // Listen to user updates
    _realtimeService.userUpdateStream.listen((update) {
      _latestUserUpdate = update;
      notifyListeners();
    });

    // Listen to address updates
    _realtimeService.addressUpdateStream.listen((update) {
      _latestAddressUpdate = update;
      notifyListeners();
    });

    // Listen to favorite updates
    _realtimeService.favoriteUpdateStream.listen((update) {
      _latestFavoriteUpdate = update;
      notifyListeners();
    });

    // Listen to notification updates
    _realtimeService.notificationUpdateStream.listen((update) {
      _latestNotificationUpdate = update;
      notifyListeners();
    });
  }

  Future<void> connectToAll() async {
    await _realtimeService.connectToAll();
    notifyListeners();
  }

  Future<void> connectToUserUpdates() async {
    await _realtimeService.connectToUserUpdates();
    notifyListeners();
  }

  Future<void> connectToAddressUpdates() async {
    await _realtimeService.connectToAddressUpdates();
    notifyListeners();
  }

  Future<void> connectToFavoriteFoodsUpdates() async {
    await _realtimeService.connectToFavoriteFoodsUpdates();
    notifyListeners();
  }

  Future<void> connectToFavoriteRestaurantsUpdates() async {
    await _realtimeService.connectToFavoriteRestaurantsUpdates();
    notifyListeners();
  }

  Future<void> connectToNotificationUpdates() async {
    await _realtimeService.connectToNotificationUpdates();
    notifyListeners();
  }

  void disconnectAll() {
    _realtimeService.disconnectAll();
    notifyListeners();
  }

  @override
  void dispose() {
    _realtimeService.dispose();
    super.dispose();
  }
}
```

---

## Usage Examples

### Example 1: Initialize in App Start

**File:** `lib/main.dart`

```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'services/realtime_service.dart';
import 'providers/realtime_provider.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Create the RealtimeService instance
    final realtimeService = RealtimeService();

    return MultiProvider(
      providers: [
        // Provide the RealtimeService
        Provider<RealtimeService>.value(value: realtimeService),

        // Provide the RealtimeProvider (wraps the service)
        ChangeNotifierProvider(
          create: (_) => RealtimeProvider(realtimeService),
        ),
      ],
      child: MaterialApp(
        title: 'DFood App',
        theme: ThemeData(
          primarySwatch: Colors.orange,
        ),
        home: const SplashScreen(),
      ),
    );
  }
}

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _initializeApp();
  }

  Future<void> _initializeApp() async {
    // Get user info (from SharedPreferences, secure storage, etc.)
    final userId = await _getUserId();
    final authToken = await _getAuthToken();

    if (userId != null) {
      // Initialize RealtimeService
      final realtimeService = context.read<RealtimeService>();
      realtimeService.initialize(
        userId: userId,
        baseUrl: 'http://your-server.com',  // Change to your server URL
        authToken: authToken,
      );

      // Connect to all endpoints
      await realtimeService.connectToAll();
    }

    // Navigate to home
    if (mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (_) => const HomeScreen()),
      );
    }
  }

  Future<String?> _getUserId() async {
    // Implement your logic to get user ID
    // Example: return await SecureStorage.read('user_id');
    return 'user123';  // Placeholder
  }

  Future<String?> _getAuthToken() async {
    // Implement your logic to get auth token
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: CircularProgressIndicator(),
      ),
    );
  }
}
```

### Example 2: Display Connection Status

**File:** `lib/widgets/connection_status_widget.dart`

```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/realtime_provider.dart';

class ConnectionStatusWidget extends StatelessWidget {
  const ConnectionStatusWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<RealtimeProvider>(
      builder: (context, realtimeProvider, child) {
        final isConnected = realtimeProvider.isConnected;
        final connectionCount = realtimeProvider.connectionCount;

        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: isConnected
                ? Colors.green.withOpacity(0.1)
                : Colors.red.withOpacity(0.1),
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: isConnected ? Colors.green : Colors.red,
              width: 1,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                isConnected ? Icons.cloud_done : Icons.cloud_off,
                size: 16,
                color: isConnected ? Colors.green : Colors.red,
              ),
              const SizedBox(width: 6),
              Text(
                isConnected
                    ? 'Live ($connectionCount)'
                    : 'Offline',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: isConnected ? Colors.green : Colors.red,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
```

### Example 3: Listen to User Profile Updates

**File:** `lib/screens/profile_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/realtime_service.dart';
import '../models/realtime_models.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  Map<String, dynamic>? _currentUser;
  String? _lastUpdate;

  @override
  void initState() {
    super.initState();
    _loadUserProfile();
    _listenToUpdates();
  }

  Future<void> _loadUserProfile() async {
    // Load user profile from API or local storage
    // setState(() { _currentUser = ...; });
  }

  void _listenToUpdates() {
    final realtimeService = context.read<RealtimeService>();

    // Listen to user updates
    realtimeService.userUpdateStream.listen((update) {
      setState(() {
        _currentUser = update.user;
        _lastUpdate = 'Updated: ${DateTime.now().toString()}';
      });

      // Show snackbar
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            'Profile updated: ${update.changes.keys.join(', ')}',
          ),
          duration: const Duration(seconds: 2),
          backgroundColor: Colors.green,
        ),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_currentUser == null) {
      return const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profile'),
        actions: [
          // Show connection status
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: const ConnectionStatusWidget(),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Profile Image
            Center(
              child: CircleAvatar(
                radius: 60,
                backgroundImage: _currentUser!['profile_image_url'] != null
                    ? NetworkImage(_currentUser!['profile_image_url'])
                    : null,
                child: _currentUser!['profile_image_url'] == null
                    ? const Icon(Icons.person, size: 60)
                    : null,
              ),
            ),
            const SizedBox(height: 24),

            // Name
            _buildInfoCard(
              'Name',
              '${_currentUser!['first_name']} ${_currentUser!['last_name']}',
            ),

            // Email
            _buildInfoCard('Email', _currentUser!['email']),

            // Phone
            if (_currentUser!['phone_number'] != null)
              _buildInfoCard('Phone', _currentUser!['phone_number']),

            // Bio
            if (_currentUser!['bio'] != null)
              _buildInfoCard('Bio', _currentUser!['bio']),

            // Last update indicator
            if (_lastUpdate != null)
              Padding(
                padding: const EdgeInsets.only(top: 16),
                child: Text(
                  _lastUpdate!,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                    fontStyle: FontStyle.italic,
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoCard(String label, String value) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[600],
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              value,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w400,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

### Example 4: Real-time Address List

**File:** `lib/screens/address_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/realtime_service.dart';
import '../models/realtime_models.dart';

class AddressScreen extends StatefulWidget {
  const AddressScreen({super.key});

  @override
  State<AddressScreen> createState() => _AddressScreenState();
}

class _AddressScreenState extends State<AddressScreen> {
  List<Map<String, dynamic>> _addresses = [];

  @override
  void initState() {
    super.initState();
    _loadAddresses();
    _listenToUpdates();
  }

  Future<void> _loadAddresses() async {
    // Load addresses from API
    // setState(() { _addresses = ...; });
  }

  void _listenToUpdates() {
    final realtimeService = context.read<RealtimeService>();

    realtimeService.addressUpdateStream.listen((update) {
      final address = update.address as Map<String, dynamic>;
      final action = update.action;

      setState(() {
        switch (action) {
          case 'add':
            // Add new address
            _addresses.add(address);
            _showSnackbar('Address added', Colors.green);
            break;

          case 'update':
            // Update existing address
            final index = _addresses.indexWhere(
              (a) => a['id'] == address['id'],
            );
            if (index != -1) {
              _addresses[index] = address;
              _showSnackbar('Address updated', Colors.blue);
            }
            break;

          case 'delete':
            // Remove address
            _addresses.removeWhere((a) => a['id'] == address['id']);
            _showSnackbar('Address deleted', Colors.orange);
            break;

          case 'set_default':
            // Update default address
            for (var addr in _addresses) {
              addr['is_default'] = addr['id'] == address['id'];
            }
            _showSnackbar('Default address updated', Colors.purple);
            break;
        }
      });
    });
  }

  void _showSnackbar(String message, Color color) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        duration: const Duration(seconds: 2),
        backgroundColor: color,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Addresses'),
        actions: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: const ConnectionStatusWidget(),
          ),
        ],
      ),
      body: _addresses.isEmpty
          ? const Center(child: Text('No addresses found'))
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _addresses.length,
              itemBuilder: (context, index) {
                final address = _addresses[index];
                final isDefault = address['is_default'] == true;

                return Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  elevation: isDefault ? 3 : 1,
                  child: ListTile(
                    leading: Icon(
                      _getIconForType(address['type']),
                      color: isDefault ? Colors.orange : Colors.grey,
                    ),
                    title: Text(
                      address['title'] ?? address['type'],
                      style: TextStyle(
                        fontWeight: isDefault
                            ? FontWeight.bold
                            : FontWeight.normal,
                      ),
                    ),
                    subtitle: Text(
                      '${address['street']}, ${address['city']}, ${address['state']}',
                    ),
                    trailing: isDefault
                        ? Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 8,
                              vertical: 4,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.orange.withOpacity(0.2),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: const Text(
                              'DEFAULT',
                              style: TextStyle(
                                fontSize: 10,
                                fontWeight: FontWeight.bold,
                                color: Colors.orange,
                              ),
                            ),
                          )
                        : null,
                  ),
                );
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Navigate to add address screen
        },
        child: const Icon(Icons.add),
      ),
    );
  }

  IconData _getIconForType(String type) {
    switch (type) {
      case 'home':
        return Icons.home;
      case 'work':
        return Icons.work;
      default:
        return Icons.location_on;
    }
  }
}
```

### Example 5: Real-time Notifications with Badge

**File:** `lib/screens/home_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/realtime_service.dart';
import '../models/realtime_models.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _unreadNotificationCount = 0;

  @override
  void initState() {
    super.initState();
    _loadUnreadCount();
    _listenToNotifications();
  }

  Future<void> _loadUnreadCount() async {
    // Load unread count from API
    // setState(() { _unreadNotificationCount = ...; });
  }

  void _listenToNotifications() {
    final realtimeService = context.read<RealtimeService>();

    realtimeService.notificationUpdateStream.listen((update) {
      final notification = update.notification as Map<String, dynamic>;
      final action = update.action;

      if (action == 'new') {
        setState(() {
          _unreadNotificationCount++;
        });

        // Show in-app notification
        _showInAppNotification(
          notification['title'],
          notification['body'],
        );
      }
    });
  }

  void _showInAppNotification(String title, String body) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              title,
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 4),
            Text(body),
          ],
        ),
        duration: const Duration(seconds: 4),
        backgroundColor: Colors.orange,
        action: SnackBarAction(
          label: 'VIEW',
          textColor: Colors.white,
          onPressed: () {
            // Navigate to notifications screen
          },
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('DFood'),
        actions: [
          // Connection status
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: const ConnectionStatusWidget(),
          ),

          // Notifications with badge
          Stack(
            children: [
              IconButton(
                icon: const Icon(Icons.notifications),
                onPressed: () {
                  // Navigate to notifications
                },
              ),
              if (_unreadNotificationCount > 0)
                Positioned(
                  right: 8,
                  top: 8,
                  child: Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      color: Colors.red,
                      shape: BoxShape.circle,
                    ),
                    constraints: const BoxConstraints(
                      minWidth: 16,
                      minHeight: 16,
                    ),
                    child: Text(
                      _unreadNotificationCount > 99
                          ? '99+'
                          : _unreadNotificationCount.toString(),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
            ],
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.restaurant, size: 100, color: Colors.orange),
            const SizedBox(height: 16),
            const Text(
              'Welcome to DFood!',
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            const Text('Your food will arrive soon'),
          ],
        ),
      ),
    );
  }
}
```

---

## Best Practices

### 1. Connect at App Launch

Connect to WebSocket once when the app starts (after login):

```dart
// In main.dart or after successful login
final realtimeService = context.read<RealtimeService>();
realtimeService.initialize(
  userId: currentUser.id,
  baseUrl: AppConfig.baseUrl,
  authToken: currentUser.token,
);
await realtimeService.connectToAll();
```

### 2. Disconnect on Logout

Always disconnect when user logs out:

```dart
Future<void> logout() async {
  final realtimeService = context.read<RealtimeService>();
  realtimeService.disconnectAll();

  // Clear user data, navigate to login, etc.
}
```

### 3. Handle App Lifecycle

Disconnect when app goes to background, reconnect when it returns:

```dart
class MyApp extends StatefulWidget {
  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    final realtimeService = context.read<RealtimeService>();

    switch (state) {
      case AppLifecycleState.resumed:
        // App came to foreground - reconnect
        realtimeService.connectToAll();
        break;
      case AppLifecycleState.paused:
        // App went to background - disconnect to save battery
        realtimeService.disconnectAll();
        break;
      default:
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(/* ... */);
  }
}
```

### 4. Use StreamBuilder for Reactive UI

```dart
class NotificationBadge extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final realtimeService = context.read<RealtimeService>();

    return StreamBuilder<NotificationUpdateData>(
      stream: realtimeService.notificationUpdateStream,
      builder: (context, snapshot) {
        if (snapshot.hasData && snapshot.data!.action == 'new') {
          // Update UI when new notification arrives
          return Badge(/* ... */);
        }
        return Icon(Icons.notifications);
      },
    );
  }
}
```

### 5. Combine with Local Cache

Don't rely solely on WebSocket - always have a local cache:

```dart
void _listenToUpdates() {
  final realtimeService = context.read<RealtimeService>();

  realtimeService.addressUpdateStream.listen((update) async {
    // Update local cache first (instant UI update)
    _updateLocalCache(update);

    // Optionally refresh from API to ensure consistency
    await _refreshFromAPI();
  });
}
```

### 6. Show Connection Status to User

Always show users if they're connected or offline:

```dart
// Use the ConnectionStatusWidget we created earlier
AppBar(
  title: Text('My Screen'),
  actions: [
    ConnectionStatusWidget(),
  ],
)
```

### 7. Error Handling

Handle WebSocket errors gracefully:

```dart
realtimeService.messageStream.listen(
  (message) {
    // Handle message
  },
  onError: (error) {
    print('WebSocket error: $error');
    _showErrorToUser('Connection issue. Retrying...');
  },
);
```

---

## Troubleshooting

### Problem 1: Connection Fails Immediately

**Symptom:** WebSocket connects then immediately disconnects.

**Solutions:**

1. **Check URL format:**
```dart
// ✓ Correct
baseUrl: 'http://192.168.1.100:8080'  // For local testing
baseUrl: 'https://api.yourapp.com'     // For production

// ✗ Wrong
baseUrl: 'http://192.168.1.100:8080/'  // Extra slash
baseUrl: 'ws://...'                     // Don't use ws:// in baseUrl
```

2. **Check authentication:**
- If backend requires auth, make sure you're sending token
- For web_socket_channel, you may need to send token in URL:
```dart
final uri = Uri.parse('$wsUrl/api/v1$endpoint?token=$_authToken');
```

### Problem 2: Not Receiving Messages

**Symptom:** Connected but no messages appear when data changes on backend.

**Checklist:**

1. **Verify connection is active:**
```dart
print('Connected: ${realtimeService.isConnected('user')}');
print('Connection count: ${realtimeService.connectionCount}');
```

2. **Check if you're listening to the stream:**
```dart
realtimeService.userUpdateStream.listen((update) {
  print('Received update: $update');  // Should print when update occurs
});
```

3. **Verify user ID matches:**
```dart
// User ID used in initialization must match user making changes
realtimeService.initialize(userId: 'user123', ...);
// Updates will only be sent to 'user123'
```

4. **Check backend logs:**
```bash
# On your Go server, you should see:
Realtime: User user123 connected
Realtime: Sent user update to user123, changes: map[first_name:John]
```

### Problem 3: App Crashes on Hot Reload

**Symptom:** App crashes when doing hot reload in development.

**Solution:** Properly dispose WebSocket connections:

```dart
@override
void dispose() {
  realtimeService.disconnectAll();
  super.dispose();
}
```

### Problem 4: High Battery Drain

**Symptom:** App consumes a lot of battery.

**Solutions:**

1. **Disconnect when app in background:**
```dart
@override
void didChangeAppLifecycleState(AppLifecycleState state) {
  if (state == AppLifecycleState.paused) {
    realtimeService.disconnectAll();
  }
}
```

2. **Connect only to needed endpoints:**
```dart
// Instead of connectToAll()
realtimeService.connectToUserUpdates();
realtimeService.connectToNotificationUpdates();
// Don't connect to everything if not needed
```

### Problem 5: JSON Parsing Errors

**Symptom:** `FormatException: Unexpected character` or similar.

**Solution:** Make sure your models match backend format:

```dart
// If backend sends snake_case
factory UserUpdateData.fromJson(Map<String, dynamic> json) {
  return UserUpdateData(
    user: json['user'],          // Not 'User'
    changes: json['changes'],    // Not 'Changes'
  );
}
```

### Problem 6: iOS Not Working

**Symptom:** Works on Android but not iOS.

**Solutions:**

1. **Add to Info.plist (for HTTP connections):**
```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsArbitraryLoads</key>
    <true/>
</dict>
```

2. **Use HTTPS in production** (iOS requires secure connections)

---

## Summary

### Quick Start Checklist

- [ ] Add `web_socket_channel` to `pubspec.yaml`
- [ ] Create `realtime_models.dart` with message types
- [ ] Create `realtime_service.dart` with connection management
- [ ] Create `realtime_provider.dart` (optional, for state management)
- [ ] Initialize service on app start/login
- [ ] Connect to needed endpoints
- [ ] Listen to streams in your UI
- [ ] Update UI when messages received
- [ ] Handle disconnection and reconnection
- [ ] Disconnect on logout
- [ ] Test with backend

### Key Points to Remember

1. **Initialize once** - Call `initialize()` after login, before using any other methods
2. **Connect strategically** - Only connect to endpoints you need
3. **Listen in initState** - Set up stream listeners when widgets are created
4. **Update UI reactively** - Use StreamBuilder or Provider for automatic UI updates
5. **Handle lifecycle** - Disconnect when app backgrounded, reconnect when resumed
6. **Show status** - Always show users if they're connected or offline
7. **Clean up** - Dispose connections and streams properly

### Performance Tips

- Only connect to endpoints you actively need
- Disconnect when app goes to background
- Use `broadcast` streams for multiple listeners
- Don't rebuild entire UI on every message - update only affected widgets
- Cache data locally, use WebSocket for updates only

---

**Need Help?**

Check the backend guide: `WEBSOCKET_IMPLEMENTATION_GUIDE.md`

For more information on web_socket_channel package:
https://pub.dev/packages/web_socket_channel

Happy coding! 🚀
