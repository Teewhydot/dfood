# dfood Setup Guide

## Environment Setup

### 1. Copy Environment File
```bash
copy .env.example .env
```

### 2. Configure SendGrid
1. Go to [SendGrid Dashboard](https://app.sendgrid.com/)
2. Navigate to Settings → API Keys
3. Create a new API key with "Mail Send" permissions
4. Copy the API key (starts with `SG.`)
5. Open `.env` file and replace `your_sendgrid_api_key_here` with your actual key:

```env
SENDGRID_API_KEY=SG.your_actual_api_key_here
```

### 3. Run the Application
```bash
go run cmd/main.go
```

## Environment Variables

The application supports these environment variables in `.env`:

- `SENDGRID_API_KEY` - Your SendGrid API key (required for emails)
- `APP_ENV` - Application environment (dev/staging/production)
- `PORT` - Server port (optional, overrides config file)
- `DB_DRIVER` - Database driver (optional, overrides config file)
- `DB_DATASOURCE` - Database connection string (optional, overrides config file)

## Testing Email

After setup, test email functionality by:
1. Starting the application: `go run cmd/main.go`
2. Register a new user via API
3. Check logs for email sending status

## Troubleshooting

### SendGrid 401 Error
- Verify your API key is correct and starts with `SG.`
- Ensure the API key has "Mail Send" permissions
- Check that `.env` file is in the project root
- Restart the application after changing `.env`

### Email Not Sending
- Check application logs for detailed error messages
- Verify SendGrid account is not suspended
- Ensure sender email is verified in SendGrid (for free accounts)