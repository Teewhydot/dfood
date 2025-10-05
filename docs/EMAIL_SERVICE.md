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

### 2. Configuration
The service uses configuration from your YAML files:

```yaml
sendgrid:
  api_key: ${SENDGRID_API_KEY}
  from_email: noreply@dfood.com
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

## Future Enhancements

Potential additions:
- Email verification emails
- Review notification emails
- Weekly digest emails
- Restaurant owner notifications
- Template management system