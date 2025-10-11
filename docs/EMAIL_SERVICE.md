# Email Service Documentation

## Overview
The email service handles sending transactional emails using SendGrid. Currently supports welcome emails and password reset emails.

## Setup

### 1. Environment Variables
Set your SendGrid API key as an environment variable:

```bash
# Windows CMD
set SENDGRID_API_KEY=your_sendgrid_api_key_here

# Windows PowerShell
$env:SENDGRID_API_KEY="your_sendgrid_api_key_here"
```

### 2. Verify Sender Email in SendGrid
**IMPORTANT**: Before sending emails, you must verify the sender email address in SendGrid:

1. Go to [SendGrid Dashboard](https://app.sendgrid.com/settings/sender_auth)
2. Click "Verify a Single Sender"
3. Enter your email details and verify
4. Update the `from_email` in your config files to match the verified email

**Common Error**: `403 Forbidden - "The from address does not match a verified Sender Identity"`
- **Solution**: Ensure the `from_email` in your config matches a verified sender in SendGrid

### 3. Configuration
The service uses configuration from your YAML files:

```yaml
sendgrid:
  api_key: ${SENDGRID_API_KEY}
  from_email: abubakarissa47722@gmail.com  # Must be verified in SendGrid
  from_name: dfood
```

## Usage

### Initialize Service
```go
emailService := service.NewEmailService(service.EmailConfig{
    APIKey:    cfg.SendGrid.APIKey,
    FromEmail: cfg.SendGrid.FromEmail,
    FromName:  cfg.SendGrid.FromName,
})
```

### Send Welcome Email
```go
err := emailService.SendWelcomeEmail("user@example.com", "John Doe")
if err != nil {
    // Handle error
}
```

### Send Password Reset Email
```go
resetToken := "generated_token_here"
err := emailService.SendPasswordResetEmail("user@example.com", "John Doe", resetToken)
if err != nil {
    // Handle error
}
```

### Send Generic Email
```go
subject := "Your Custom Subject"
plainText := "Plain text version of your email"
htmlContent := "<h1>HTML version</h1><p>Your custom HTML content</p>"

err := emailService.SendEmail("user@example.com", "John Doe", subject, plainText, htmlContent)
if err != nil {
    // Handle error
}
```

## Integration with Auth Service

The email service is automatically integrated with the auth service:

- **Registration**: Sends welcome email after successful user registration
- **Password Reset**: Can be used in password reset flow (you'll need to implement the reset token generation)

## Email Templates

Templates are now organized in a separate file (`internal/service/email_templates.go`) for better maintainability.

### Available Templates

#### Welcome Email
- **Subject**: "Welcome to dfood!"
- **Content**: Welcome message with app features
- **Format**: Both HTML and plain text

#### Password Reset Email
- **Subject**: "Reset Your dfood Password"
- **Content**: Reset link with 1-hour expiration notice
- **Format**: Both HTML and plain text
- **Security**: Includes warning about ignoring if not requested

#### Review Notification Email
- **Subject**: "New Review for [Restaurant Name]"
- **Content**: Notifies restaurant owners of new reviews
- **Format**: Both HTML and plain text

#### Order Confirmation Email
- **Subject**: "Order Confirmation - [Order ID]"
- **Content**: Order details and confirmation
- **Format**: Both HTML and plain text

### Adding New Templates

To add new email templates:

1. Add the template function to `internal/service/email_templates.go`
2. Add the corresponding method to the `EmailService` interface
3. Implement the method in `emailService` struct

## Error Handling

The service returns errors for:
- SendGrid API failures
- Network issues
- Invalid API keys
- Rate limiting

Errors are wrapped with context for better debugging.

## Best Practices

1. **Async Sending**: Welcome emails are sent asynchronously to avoid blocking registration
2. **Error Logging**: Log email failures but don't fail the main operation
3. **Template Updates**: Update email templates in the service code
4. **Testing**: Use SendGrid's sandbox mode for testing

## Troubleshooting

### Error: 403 Forbidden - Sender Identity Not Verified
**Problem**: `The from address does not match a verified Sender Identity`

**Solution**:
1. Verify your sender email at https://app.sendgrid.com/settings/sender_auth
2. Ensure config `from_email` matches the verified email exactly
3. Restart your application after config changes

### Error: 401 Unauthorized
**Problem**: Invalid SendGrid API key

**Solution**:
1. Check `SENDGRID_API_KEY` environment variable is set correctly
2. Verify API key hasn't expired in SendGrid dashboard
3. Ensure no extra spaces in the API key

### Email Not Received
**Check**:
1. Verify sender email is verified in SendGrid
2. Check spam/junk folder
3. Review SendGrid activity logs at https://app.sendgrid.com/email_activity
4. Ensure API key has sending permissions

### Testing Email Sending
```bash
# Restart app after config changes
go run cmd/main.go

# Monitor logs for SendGrid responses
# Look for: "Email sent successfully - Status: 202"
```

## Future Enhancements

Potential additions:
- Email verification emails
- Review notification emails
- Weekly digest emails
- Restaurant owner notifications
- Template management system