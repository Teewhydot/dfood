# Postman Import Guide

## Quick Start (3 Steps)

### 1. Import Collection
1. Open Postman
2. Click **Import** button (top left)
3. Drag and drop `postman-collection.json` or click **Upload Files**
4. Click **Import**

### 2. Create Environment
1. Click **Environments** (left sidebar)
2. Click **+** to create new environment
3. Name it: `DFood Local`
4. Add these variables:

| Variable | Initial Value | Current Value |
|----------|--------------|---------------|
| `base_url` | `http://localhost:8080/api/v1` | `http://localhost:8080/api/v1` |
| `access_token` | (leave empty) | (leave empty) |
| `refresh_token` | (leave empty) | (leave empty) |

5. Click **Save**
6. Select `DFood Local` from environment dropdown (top right)

### 3. Start Testing
1. Start your backend server: `go run cmd/main.go`
2. Go to **Auth Complete** folder
3. Run **Register a new user** request
4. Run **Login to get tokens** request
   - Tokens will auto-populate in environment variables
5. All other requests will now work with the bearer token

## What's Included

**84 requests across 15 folders:**
- Addresses
- Auth Complete (register, login, logout, password reset)
- Auth Middleware Test
- Crud Endpoints
- Favorites
- Foods (search, categories, popular)
- Notifications
- Order Test
- Password Reset Flow
- Payments
- Rate Limit Test
- Restaurants (search, nearby, menu)
- Uploads
- Users (profile, FCM tokens)
- Workflow

## Auto Token Capture

Login and register requests automatically save tokens to your environment variables. No manual copy/paste needed.

## Tips

- Use `{{base_url}}` instead of full URL if you want to switch environments
- Clone `DFood Local` to create `DFood Staging` or `DFood Production` environments
- Use Postman Collections Runner to run multiple requests in sequence
- Save example responses for documentation

## Troubleshooting

**Requests failing with 401?**
- Run login request again to refresh tokens
- Check environment is selected (top right dropdown)

**Variables not showing?**
- Make sure environment is selected
- Check variable names match exactly (case-sensitive)

**Can't see folders?**
- Collection imported as flat list? Re-import the file
- Folders should appear automatically

## Re-generating Collection

If you update your `.http` files, regenerate the collection:

```bash
python3 convert_to_postman.py
```

Then re-import in Postman (it will update existing collection).
